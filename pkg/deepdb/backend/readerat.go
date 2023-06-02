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

package backend

import (
	"context"
	"io"

	"github.com/intergral/deep/pkg/util"
)

// ContextReader is an io.ReaderAt interface that passes context.  It is used to simplify access to backend objects
// and abstract away the name/meta and other details so that the data can be accessed directly and simply
type ContextReader interface {
	ReadAt(ctx context.Context, p []byte, off int64) (int, error)
	ReadAll(ctx context.Context) ([]byte, error)

	// Return an io.Reader representing the underlying. May not be supported by all implementations
	Reader() (io.Reader, error)
}

// backendReader is a shim that allows a backend.Reader to be used as a ContextReader
type backendReader struct {
	meta        *BlockMeta
	name        string
	r           Reader
	shouldCache bool
}

// NewContextReader creates a ReaderAt for the given BlockMeta
func NewContextReader(meta *BlockMeta, name string, r Reader, shouldCache bool) ContextReader {
	return &backendReader{
		meta:        meta,
		name:        name,
		r:           r,
		shouldCache: shouldCache,
	}
}

// ReadAt implements ContextReader
func (b *backendReader) ReadAt(ctx context.Context, p []byte, off int64) (int, error) {
	err := b.r.ReadRange(ctx, b.name, b.meta.BlockID, b.meta.TenantID, uint64(off), p, false)
	return len(p), err
}

// ReadAll implements ContextReader
func (b *backendReader) ReadAll(ctx context.Context) ([]byte, error) {
	return b.r.Read(ctx, b.name, b.meta.BlockID, b.meta.TenantID, b.shouldCache)
}

// Reader implements ContextReader
func (b *backendReader) Reader() (io.Reader, error) {
	return nil, util.ErrUnsupported
}

// AllReader is an interface that supports both io.Reader and io.ReaderAt methods
type AllReader interface {
	io.Reader
	io.ReaderAt
}

// allReader wraps an AllReader and implements backend.ContextReader
type allReader struct {
	r AllReader
}

// NewContextReaderWithAllReader wraps a normal ReaderAt and drops the context
func NewContextReaderWithAllReader(r AllReader) ContextReader {
	return &allReader{
		r: r,
	}
}

// ReadAt implements ContextReader
func (r *allReader) ReadAt(ctx context.Context, p []byte, off int64) (int, error) {
	return r.r.ReadAt(p, off)
}

// ReadAll implements ContextReader
func (r *allReader) ReadAll(ctx context.Context) ([]byte, error) {
	return io.ReadAll(r.r)
}

// Reader implements ContextReader
func (r *allReader) Reader() (io.Reader, error) {
	return r.r, nil
}
