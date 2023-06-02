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

package registry

import (
	"math"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestTestRegistry_counter(t *testing.T) {
	testRegistry := NewTestRegistry()

	counter := testRegistry.NewCounter("counter", []string{"foo", "bar"})

	labelValues := newLabelValues([]string{"foo-value", "bar-value"})
	counter.Inc(labelValues, 1.0)
	counter.Inc(labelValues, 2.0)
	counter.Inc(labelValues, 1.5)

	lbls := labels.FromMap(map[string]string{
		"foo": "foo-value",
		"bar": "bar-value",
	})
	assert.Equal(t, 4.5, testRegistry.Query("counter", lbls))
}

func TestTestRegistry_histogram(t *testing.T) {
	testRegistry := NewTestRegistry()

	histogram := testRegistry.NewHistogram("histogram", []string{"foo", "bar"}, []float64{1.0, 2.0})

	labelValues := newLabelValues([]string{"foo-value", "bar-value"})
	histogram.ObserveWithExemplar(labelValues, 1.0, "", 1.0)
	histogram.ObserveWithExemplar(labelValues, 2.0, "", 1.0)
	histogram.ObserveWithExemplar(labelValues, 2.5, "", 1.0)

	lbls := labels.FromMap(map[string]string{
		"foo": "foo-value",
		"bar": "bar-value",
	})
	assert.Equal(t, 1.0, testRegistry.Query("histogram_bucket", withLe(lbls, 1.0)))
	assert.Equal(t, 2.0, testRegistry.Query("histogram_bucket", withLe(lbls, 2.0)))
	assert.Equal(t, 3.0, testRegistry.Query("histogram_bucket", withLe(lbls, math.Inf(1))))
	assert.Equal(t, 3.0, testRegistry.Query("histogram_count", lbls))
	assert.Equal(t, 5.5, testRegistry.Query("histogram_sum", lbls))
}
