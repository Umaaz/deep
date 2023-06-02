/*
 * Copyright (C) 2023  Intergral GmbH
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package vparquet

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/google/uuid"
	"go.uber.org/atomic"

	"github.com/intergral/deep/pkg/deepdb/backend"
	"github.com/intergral/deep/pkg/deepdb/encoding/common"
)

// This stack of readers is used to bridge the gap between the backend.Reader and the parquet.File.
//  each fulfills a different role.
// backend.Reader <- BackendReaderAt <- io.BufferedReaderAt <- parquetOptimizedReaderAt <- cachedReaderAt <- parquet.File
//                                \                                                         /
//                                  <------------------------------------------------------

// BackendReaderAt is used to track backend requests and present a io.ReaderAt interface backed
// by a backend.Reader
type BackendReaderAt struct {
	ctx      context.Context
	r        backend.Reader
	name     string
	blockID  uuid.UUID
	tenantID string

	TotalBytesRead atomic.Uint64
}

var _ io.ReaderAt = (*BackendReaderAt)(nil)

func NewBackendReaderAt(ctx context.Context, r backend.Reader, name string, blockID uuid.UUID, tenantID string) *BackendReaderAt {
	return &BackendReaderAt{ctx, r, name, blockID, tenantID, atomic.Uint64{}}
}

func (b *BackendReaderAt) ReadAt(p []byte, off int64) (int, error) {
	b.TotalBytesRead.Add(uint64(len(p)))
	err := b.r.ReadRange(b.ctx, b.name, b.blockID, b.tenantID, uint64(off), p, false)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

func (b *BackendReaderAt) ReadAtWithCache(p []byte, off int64) (int, error) {
	err := b.r.ReadRange(b.ctx, b.name, b.blockID, b.tenantID, uint64(off), p, true)
	if err != nil {
		return 0, err
	}
	return len(p), err
}

// parquetOptimizedReaderAt is used to cheat a few parquet calls. By default when opening a
// file parquet always requests the magic number and then the footer length. We can save
// both of these calls from going to the backend.
type parquetOptimizedReaderAt struct {
	r          io.ReaderAt
	readerSize int64
	footerSize uint32
}

var _ io.ReaderAt = (*parquetOptimizedReaderAt)(nil)

func newParquetOptimizedReaderAt(r io.ReaderAt, size int64, footerSize uint32) *parquetOptimizedReaderAt {
	return &parquetOptimizedReaderAt{r, size, footerSize}
}

func (r *parquetOptimizedReaderAt) ReadAt(p []byte, off int64) (int, error) {
	if len(p) == 4 && off == 0 {
		// Magic header
		return copy(p, []byte("PAR1")), nil
	}

	if len(p) == 8 && off == r.readerSize-8 && r.footerSize > 0 /* not present in previous block metas */ {
		// Magic footer
		binary.LittleEndian.PutUint32(p, r.footerSize)
		copy(p[4:8], []byte("PAR1"))
		return 8, nil
	}

	return r.r.ReadAt(p, off)
}

// cachedReaderAt is used to route specific reads to the caching layer. this must be passed directly into
// the parquet.File so thet Set*Section() methods get called.
type cachedReaderAt struct {
	r             io.ReaderAt
	br            *BackendReaderAt
	cacheControl  common.CacheControl
	cachedObjects map[int64]int64 // storing offsets and length of objects we want to cache
}

var _ io.ReaderAt = (*cachedReaderAt)(nil)

func newCachedReaderAt(br io.ReaderAt, rr *BackendReaderAt, cc common.CacheControl) *cachedReaderAt {
	return &cachedReaderAt{br, rr, cc, map[int64]int64{}}
}

// called by parquet-go in OpenFile() to set offset and length of footer section
func (r *cachedReaderAt) SetFooterSection(offset, length int64) {
	if r.cacheControl.Footer {
		r.cachedObjects[offset] = length
	}
}

// called by parquet-go in OpenFile() to set offset and length of column indexes
func (r *cachedReaderAt) SetColumnIndexSection(offset, length int64) {
	if r.cacheControl.ColumnIndex {
		r.cachedObjects[offset] = length
	}
}

// called by parquet-go in OpenFile() to set offset and length of offset index section
func (r *cachedReaderAt) SetOffsetIndexSection(offset, length int64) {
	if r.cacheControl.OffsetIndex {
		r.cachedObjects[offset] = length
	}
}

func (r *cachedReaderAt) ReadAt(p []byte, off int64) (int, error) {
	// check if the offset and length is stored as a special object
	if r.cachedObjects[off] == int64(len(p)) {
		return r.br.ReadAtWithCache(p, off)
	}

	return r.r.ReadAt(p, off)
}
