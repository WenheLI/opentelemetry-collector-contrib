// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 * Utility functions to extract metadata from OTEL metrics
 */

package metadataexporter

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func _extractNumberDataPoints(dataPoints pmetric.NumberDataPointSlice, metricMetadata *MetricMetadataPoint) {
	var lastetTime pcommon.Timestamp
	for i := 0; i < dataPoints.Len(); i++ {
		dp := dataPoints.At(i)
		if dp.Timestamp() > lastetTime {
			lastetTime = dp.Timestamp()
		}
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			// TODO handle panic
			metricMetadata.Dimensions[k] = v.Type().String()
			return true
		})
	}
	metricMetadata.LastPublishtime = lastetTime.AsTime().UnixMilli()

}

func _extractHistogramDataPoints(dataPoints pmetric.HistogramDataPointSlice, metricMetadata *MetricMetadataPoint) {
	var lastetTime pcommon.Timestamp
	for i := 0; i < dataPoints.Len(); i++ {
		dp := dataPoints.At(i)
		if dp.Timestamp() > lastetTime {
			lastetTime = dp.Timestamp()
		}
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			// TODO handle panic
			metricMetadata.Dimensions[k] = v.Type().String()
			return true
		})
	}
	metricMetadata.LastPublishtime = lastetTime.AsTime().UnixMilli()
}

func _extractExponentialHistogramDataPoints(dataPoints pmetric.ExponentialHistogramDataPointSlice, metricMetadata *MetricMetadataPoint) {
	var lastetTime pcommon.Timestamp
	for i := 0; i < dataPoints.Len(); i++ {
		dp := dataPoints.At(i)
		if dp.Timestamp() > lastetTime {
			lastetTime = dp.Timestamp()
		}
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			// TODO handle panic
			metricMetadata.Dimensions[k] = v.Type().String()
			return true
		})
	}
	metricMetadata.LastPublishtime = lastetTime.AsTime().UnixMilli()
}

func _extractSummaryDataPoints(dataPoints pmetric.SummaryDataPointSlice, metricMetadata *MetricMetadataPoint) {
	var lastetTime pcommon.Timestamp
	for i := 0; i < dataPoints.Len(); i++ {
		dp := dataPoints.At(i)
		if dp.Timestamp() > lastetTime {
			lastetTime = dp.Timestamp()
		}
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			metricMetadata.Dimensions[k] = v.Type().String()
			return true
		})
	}
	metricMetadata.LastPublishtime = lastetTime.AsTime().UnixMilli()
}

func extractMetric(scopeMetrics pmetric.ScopeMetricsSlice) []MetricMetadataPoint {
	metadataList := make([]MetricMetadataPoint, 0)
	metadataNameHash := make(map[string]bool)

	for i := 0; i < scopeMetrics.Len(); i++ {
		scopeMetrics := scopeMetrics.At(i)
		isSLI := scopeMetrics.IsSLI()
		sliDetailMap := scopeMetrics.SliDetail()
		sliDetail := make(map[string]string)
		sliDetailMap.Range((func(k string, v pcommon.Value) bool {
			sliDetail[k] = v.StringVal()
			return true
		}))
		metrics := scopeMetrics.Metrics()
		for j := 0; j < metrics.Len(); j++ {
			metric := metrics.At(j)

			metadataName := metric.Name()
			if metadataNameHash[metadataName] {
				continue // skip duplicate metric name
			}
			metadataDescription := metric.Description()

			metricMetadata := MetricMetadataPoint{
				Name:        metadataName,
				Description: metadataDescription,
				Dimensions:  make(map[string]string),
				IsSLI:       isSLI,
				SLIDetail:   sliDetail,
			}

			switch metric.DataType() {
			case pmetric.MetricDataTypeGauge:
				_extractNumberDataPoints(metric.Gauge().DataPoints(), &metricMetadata)
			case pmetric.MetricDataTypeSum:
				_extractNumberDataPoints(metric.Sum().DataPoints(), &metricMetadata)
			case pmetric.MetricDataTypeHistogram:
				_extractHistogramDataPoints(metric.Histogram().DataPoints(), &metricMetadata)
			case pmetric.MetricDataTypeExponentialHistogram:
				_extractExponentialHistogramDataPoints(metric.ExponentialHistogram().DataPoints(), &metricMetadata)
			case pmetric.MetricDataTypeSummary:
				_extractSummaryDataPoints(metric.Summary().DataPoints(), &metricMetadata)
			case pmetric.MetricDataTypeNone:
				continue
			}
			metadataList = append(metadataList, metricMetadata)
			metadataNameHash[metadataName] = true
		}
	}

	return metadataList
}
