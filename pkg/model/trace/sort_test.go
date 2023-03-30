package trace

import (
	"testing"

	"github.com/intergral/deep/pkg/deeppb"
	v1 "github.com/intergral/deep/pkg/deeppb/trace/v1"
	"github.com/stretchr/testify/assert"
)

func TestSortTrace(t *testing.T) {
	tests := []struct {
		input    *deeppb.Trace
		expected *deeppb.Trace
	}{
		{
			input:    &deeppb.Trace{},
			expected: &deeppb.Trace{},
		},

		{
			input: &deeppb.Trace{
				Batches: []*v1.ResourceSpans{
					{
						ScopeSpans: []*v1.ScopeSpans{
							{
								Spans: []*v1.Span{
									{
										StartTimeUnixNano: 2,
									},
								},
							},
						},
					},
					{
						ScopeSpans: []*v1.ScopeSpans{
							{
								Spans: []*v1.Span{
									{
										StartTimeUnixNano: 1,
									},
								},
							},
						},
					},
				},
			},
			expected: &deeppb.Trace{
				Batches: []*v1.ResourceSpans{
					{
						ScopeSpans: []*v1.ScopeSpans{
							{
								Spans: []*v1.Span{
									{
										StartTimeUnixNano: 1,
									},
								},
							},
						},
					},
					{
						ScopeSpans: []*v1.ScopeSpans{
							{
								Spans: []*v1.Span{
									{
										StartTimeUnixNano: 2,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		SortTrace(tt.input)

		assert.Equal(t, tt.expected, tt.input)
	}
}
