package converter

import (
	"strings"

	"go.opentelemetry.io/collector/model/pdata"
)

func containsMetricWithPrefix(metricSlice pdata.MetricSlice, prefix string) bool {
	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(i)

		if strings.HasPrefix(metric.Name(), prefix) {
			return true
		}
	}

	return false
}

func containsAttributes(attributeMap pdata.AttributeMap, attributes ...string) bool {
	for i := 0; i < len(attributes); i++ {
		_, ex := attributeMap.Get(attributes[i])

		if !ex {
			return false
		}
	}

	return true
}
