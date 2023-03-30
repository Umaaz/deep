package util

import (
	v1_common "github.com/intergral/deep/pkg/deeppb/common/v1"
	v1 "github.com/intergral/deep/pkg/deeppb/trace/v1"
	deep_util "github.com/intergral/deep/pkg/util"
	semconv "go.opentelemetry.io/collector/semconv/v1.9.0"
)

func FindServiceName(attributes []*v1_common.KeyValue) (string, bool) {
	return FindAttributeValue(semconv.AttributeServiceName, attributes)
}

func FindAttributeValue(key string, attributes ...[]*v1_common.KeyValue) (string, bool) {
	for _, attrs := range attributes {
		for _, kv := range attrs {
			if key == kv.Key {
				return deep_util.StringifyAnyValue(kv.Value), true
			}
		}
	}
	return "", false
}

func GetSpanMultiplier(ratioKey string, span *v1.Span) float64 {
	spanMultiplier := 1.0
	if ratioKey != "" {
		for _, kv := range span.Attributes {
			if kv.Key == ratioKey {
				v := kv.Value.GetDoubleValue()
				if v > 0 {
					spanMultiplier = 1.0 / v
				}
			}
		}
	}
	return spanMultiplier
}
