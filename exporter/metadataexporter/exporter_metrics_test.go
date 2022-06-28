// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadataexporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type DataPointTestCase struct {
	timestamp  int64
	attributes map[string]interface{}
}

type MetricExporterTestCase struct {
	dps      []DataPointTestCase
	expected MetricMetadataPoint
}

var metricExporterTestCases = []MetricExporterTestCase{
	{
		dps: []DataPointTestCase{{
			timestamp: 1,
			attributes: map[string]interface{}{
				"dim1": "b",
			},
		}, {
			timestamp: 2,
			attributes: map[string]interface{}{
				"dim1": "c",
			},
		}},
		expected: MetricMetadataPoint{
			LastPublishtime: pcommon.NewTimestampFromTime(time.Unix(2, 0)).String(),
			Dimensions: map[string]string{
				"dim1": "STRING",
			},
		},
	},
	{
		dps: []DataPointTestCase{
			{
				timestamp: 1,
				attributes: map[string]interface{}{
					"dim1": "b",
					"dim2": "b",
					"dim3": "b",
				},
			}, {
				timestamp: 2,
				attributes: map[string]interface{}{
					"dim1": "c",
				},
			},
		},
		expected: MetricMetadataPoint{
			LastPublishtime: pcommon.NewTimestampFromTime(time.Unix(2, 0)).String(),
			Dimensions: map[string]string{
				"dim1": "STRING",
				"dim2": "STRING",
				"dim3": "STRING",
			},
		},
	},
}

func TestExtractNumberDataPoints(t *testing.T) {
	for _, testCase := range metricExporterTestCases {
		numberDPs := pmetric.NewNumberDataPointSlice()
		for _, dp := range testCase.dps {
			currentDP := numberDPs.AppendEmpty()
			currentDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(dp.timestamp, 0)))
			for k, v := range dp.attributes {
				currentDP.Attributes().Insert(k, pcommon.NewValueString(v.(string)))
			}
		}

		result := MetricMetadataPoint{}
		result.Dimensions = make(map[string]string)

		err := _extractNumberDataPoints(numberDPs, &result)

		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestExtractHistogramDataPoints(t *testing.T) {
	for _, testCase := range metricExporterTestCases {
		histogramDPs := pmetric.NewHistogramDataPointSlice()
		for _, dp := range testCase.dps {
			currentDP := histogramDPs.AppendEmpty()
			currentDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(dp.timestamp, 0)))
			for k, v := range dp.attributes {
				currentDP.Attributes().Insert(k, pcommon.NewValueString(v.(string)))
			}
		}

		result := MetricMetadataPoint{}
		result.Dimensions = make(map[string]string)

		err := _extractHistogramDataPoints(histogramDPs, &result)

		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestExponentialHistogramDataPoints(t *testing.T) {
	for _, testCase := range metricExporterTestCases {
		expoDPs := pmetric.NewExponentialHistogramDataPointSlice()
		for _, dp := range testCase.dps {
			currentDP := expoDPs.AppendEmpty()
			currentDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(dp.timestamp, 0)))
			for k, v := range dp.attributes {
				currentDP.Attributes().Insert(k, pcommon.NewValueString(v.(string)))
			}
		}

		result := MetricMetadataPoint{}
		result.Dimensions = make(map[string]string)

		err := _extractExponentialHistogramDataPoints(expoDPs, &result)

		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestExtractSummaryDataPoints(t *testing.T) {
	for _, testCase := range metricExporterTestCases {
		summaryDPs := pmetric.NewSummaryDataPointSlice()
		for _, dp := range testCase.dps {
			currentDP := summaryDPs.AppendEmpty()
			currentDP.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(dp.timestamp, 0)))
			for k, v := range dp.attributes {
				currentDP.Attributes().Insert(k, pcommon.NewValueString(v.(string)))
			}
		}

		result := MetricMetadataPoint{}
		result.Dimensions = make(map[string]string)

		err := _extractSummaryDataPoints(summaryDPs, &result)

		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

type ExtractMetricTestData struct {
	name        []string
	description []string
	dataType    []pmetric.MetricDataType
}

type ExtractMetricTestResult struct {
	name        []string
	description []string
	length      int
}

type ExtractMetricTestCase struct {
	data   ExtractMetricTestData
	result ExtractMetricTestResult
}

var extractMetricTestCases = []ExtractMetricTestCase{
	{
		data: ExtractMetricTestData{
			name:        []string{"a", "b", "c"},
			description: []string{"description_a", "description_b", "description_c"},
			dataType:    []pmetric.MetricDataType{pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram},
		},
		result: ExtractMetricTestResult{
			name:        []string{"a", "b", "c"},
			description: []string{"description_a", "description_b", "description_c"},
			length:      3,
		},
	}, {
		data: ExtractMetricTestData{
			name:        []string{"a", "b", "c"},
			description: []string{"description_a", "description_b", "description_c"},
			dataType:    []pmetric.MetricDataType{pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeNone},
		},
		result: ExtractMetricTestResult{
			name:        []string{"a", "b"},
			description: []string{"description_a", "description_b"},
			length:      2,
		},
	}, {
		data: ExtractMetricTestData{
			name:        []string{"a", "a", "c"},
			description: []string{"description_a", "description_b", "description_c"},
			dataType:    []pmetric.MetricDataType{pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram},
		},
		result: ExtractMetricTestResult{
			name:        []string{"a", "c"},
			description: []string{"description_a", "description_c"},
			length:      2,
		},
	}, {
		data: ExtractMetricTestData{
			name:        []string{"a", "b", "c"},
			description: []string{"description_a", "description_a", "description_c"},
			dataType:    []pmetric.MetricDataType{pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeHistogram},
		},
		result: ExtractMetricTestResult{
			name:        []string{"a", "b", "c"},
			description: []string{"description_a", "description_a", "description_c"},
			length:      3,
		},
	}, {
		data: ExtractMetricTestData{
			name:        []string{"a", "b", "c", "d", "e"},
			description: []string{"description_a", "description_a", "description_c", "description_d", "description_e"},
			dataType:    []pmetric.MetricDataType{pmetric.MetricDataTypeHistogram, pmetric.MetricDataTypeExponentialHistogram, pmetric.MetricDataTypeGauge, pmetric.MetricDataTypeSummary, pmetric.MetricDataTypeSum},
		},
		result: ExtractMetricTestResult{
			name:        []string{"a", "b", "c", "d", "e"},
			description: []string{"description_a", "description_a", "description_c", "description_d", "description_e"},
			length:      5,
		},
	},
}

func TestExtractMetric(t *testing.T) {
	for _, testCase := range extractMetricTestCases {
		scopeMetrics := pmetric.NewScopeMetricsSlice()
		metrics := scopeMetrics.AppendEmpty().Metrics()
		for i, name := range testCase.data.name {
			metric := metrics.AppendEmpty()
			metric.SetName(name)
			metric.SetDescription(testCase.data.description[i])
			metric.SetDataType(testCase.data.dataType[i])
		}

		result := extractMetric(scopeMetrics)
		assert.Equal(t, testCase.result.length, len(result))
		for i, name := range testCase.result.name {
			assert.Equal(t, name, result[i].Name)
			assert.Equal(t, testCase.result.description[i], result[i].Description)
		}
	}
}
