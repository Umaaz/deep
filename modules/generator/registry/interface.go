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

// Registry is a metrics store.
type Registry interface {
	NewLabelValues(values []string) *LabelValues
	NewCounter(name string, labels []string) Counter
	NewHistogram(name string, labels []string, buckets []float64) Histogram
}

// Counter
// https://prometheus.io/docs/concepts/metric_types/#counter
type Counter interface {
	Inc(values *LabelValues, value float64)
}

// Histogram
// https://prometheus.io/docs/concepts/metric_types/#histogram
type Histogram interface {
	// ObserveWithExemplar observes a datapoint with the given values. traceID will be added as exemplar.
	ObserveWithExemplar(values *LabelValues, value float64, traceID string, multiplier float64)
}

// LabelValues is a wrapper around a slice of label values. It has the ability to cache the hash of
// the label values. When passing the same label values to multiple metrics, create LabelValues once
// and pass it to all of them.
type LabelValues struct {
	values []string
	hash   uint64
}

func newLabelValues(values []string) *LabelValues {
	return &LabelValues{
		values: values,
		hash:   0,
	}
}

func newLabelValuesWithMax(values []string, maxLengthLabelValue int) *LabelValues {
	truncateLength(values, maxLengthLabelValue)
	return newLabelValues(values)
}

func (l *LabelValues) getValues() []string {
	if l == nil {
		return nil
	}
	return l.values
}

func (l *LabelValues) getValuesCopy() []string {
	labelValuesCopy := make([]string, len(l.getValues()))
	copy(labelValuesCopy, l.getValues())
	return labelValuesCopy
}

func (l *LabelValues) getHash() uint64 {
	if l == nil {
		return 0
	}
	if l.hash != 0 {
		return l.hash
	}
	l.hash = hashLabelValues(l.values)
	return l.hash
}
