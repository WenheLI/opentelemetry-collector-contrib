package azureblobexporter

import (
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var HistogramTestCase = map[interface{}]interface{}{
	"inputs": []map[interface{}]interface{}{
		{
			"start_timestamp": time.Unix(0, 0),
			"timestamp":       time.Unix(1, 0),
			"bucket_counts":   []uint64{1, 2, 3},
			"explicit_bounds": []float64{0, 1, 2, 3},
			"attributes": map[interface{}]interface{}{
				"key":   "key1",
				"value": "value1",
			},
		}, {
			"start_timestamp": time.Unix(1, 0),
			"timestamp":       time.Unix(2, 0),
			"bucket_counts":   []uint64{4, 5, 6},
			"explicit_bounds": []float64{3, 4, 5, 6},
			"attributes": map[interface{}]interface{}{
				"key":   "key2",
				"value": "value2",
			},
		},
	},
	"output": []MetricParquetStruct{
		{
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key1",
				"value": "value1",
			},
			TimeSeriesValues: map[string][]float64{
				"counts": {1, 2, 3},
				"bounds": {0, 1, 2, 3},
			},
			StartTimestamp: time.Unix(0, 0).UnixMilli(),
			EndTimestamp:   time.Unix(1, 0).UnixMilli(),
			TypeName:       "histogram",
		}, {
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key2",
				"value": "value2",
			},
			TimeSeriesValues: map[string][]float64{
				"counts": {4, 5, 6},
				"bounds": {3, 4, 5, 6},
			},
			StartTimestamp: time.Unix(1, 0).UnixMilli(),
			EndTimestamp:   time.Unix(2, 0).UnixMilli(),
			TypeName:       "histogram",
		},
	},
}

func makeTestcaseFromHistogram(testcase map[interface{}]interface{}) (pmetric.HistogramDataPointSlice, []MetricParquetStruct) {
	dps := pmetric.NewHistogramDataPointSlice()
	for _, input := range testcase["inputs"].([]map[interface{}]interface{}) {
		dp := dps.AppendEmpty()
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(input["start_timestamp"].(time.Time)))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(input["timestamp"].(time.Time)))
		dp.SetBucketCounts(pcommon.NewImmutableUInt64Slice(input["bucket_counts"].([]uint64)))
		dp.SetExplicitBounds(pcommon.NewImmutableFloat64Slice(input["explicit_bounds"].([]float64)))
		attrs := dp.Attributes()
		for k, v := range input["attributes"].(map[interface{}]interface{}) {
			attrs.Insert(k.(string), pcommon.NewValueString(v.(string)))
		}
	}
	return dps, testcase["output"].([]MetricParquetStruct)
}

func TestProcessHistogramDataPoints(t *testing.T) {
	raw, expected := makeTestcaseFromHistogram(HistogramTestCase)
	name := "test"
	dps := processHistogramDataPoints(raw, &name)
	assert.Equal(t, expected, dps)
}

var NumberTestCase = map[interface{}]interface{}{
	"inputs": []map[interface{}]interface{}{
		{
			"start_timestamp": time.Unix(0, 0),
			"timestamp":       time.Unix(1, 0),
			"doubleValue":     1.0,
			"intValue":        nil,
			"attributes": map[interface{}]interface{}{
				"key":   "key1",
				"value": "value1",
			},
		}, {
			"start_timestamp": time.Unix(1, 0),
			"timestamp":       time.Unix(2, 0),
			"doubleValue":     nil,
			"intValue":        2,
			"attributes": map[interface{}]interface{}{
				"key":   "key2",
				"value": "value2",
			},
		},
	},
	"output": []MetricParquetStruct{
		{
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key1",
				"value": "value1",
			},
			TimeSeriesValues: map[string][]float64{
				"value": {1.0},
			},
			StartTimestamp: time.Unix(0, 0).UnixMilli(),
			EndTimestamp:   time.Unix(1, 0).UnixMilli(),
			TypeName:       "number",
		}, {
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key2",
				"value": "value2",
			},
			TimeSeriesValues: map[string][]float64{
				"value": {2.0},
			},
			StartTimestamp: time.Unix(1, 0).UnixMilli(),
			EndTimestamp:   time.Unix(2, 0).UnixMilli(),
			TypeName:       "number",
		},
	},
}

func makeTestcaseFromNumber(testcase map[interface{}]interface{}) (pmetric.NumberDataPointSlice, []MetricParquetStruct) {
	dps := pmetric.NewNumberDataPointSlice()
	for _, input := range testcase["inputs"].([]map[interface{}]interface{}) {
		dp := dps.AppendEmpty()
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(input["start_timestamp"].(time.Time)))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(input["timestamp"].(time.Time)))
		if input["doubleValue"] != nil {
			dp.SetDoubleVal(input["doubleValue"].(float64))
		} else {
			dp.SetIntVal(int64(input["intValue"].(int)))
		}
		attrs := dp.Attributes()
		for k, v := range input["attributes"].(map[interface{}]interface{}) {
			attrs.Insert(k.(string), pcommon.NewValueString(v.(string)))
		}
	}
	return dps, testcase["output"].([]MetricParquetStruct)
}

func TestProcessNumberDataPoints(t *testing.T) {
	raw, expected := makeTestcaseFromNumber(NumberTestCase)
	name := "test"
	dps := processNumberDataPoints(raw, &name)
	assert.Equal(t, expected, dps)
}

var ExpHistoTestcase = map[interface{}]interface{}{
	"inputs": []map[interface{}]interface{}{
		{
			"start_timestamp": time.Unix(0, 0),
			"timestamp":       time.Unix(1, 0),
			"count":           4,
			"sum":             10,
			"attributes": map[interface{}]interface{}{
				"key":   "key1",
				"value": "value1",
			},
		}, {
			"start_timestamp": time.Unix(1, 0),
			"timestamp":       time.Unix(2, 0),
			"count":           2,
			"sum":             20,
			"attributes": map[interface{}]interface{}{
				"key":   "key2",
				"value": "value2",
			},
		},
	},
	"output": []MetricParquetStruct{
		{
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key1",
				"value": "value1",
			},
			TimeSeriesValues: map[string][]float64{
				"count": {4.0},
				"sum":   {10.0},
			},
			StartTimestamp: time.Unix(0, 0).UnixMilli(),
			EndTimestamp:   time.Unix(1, 0).UnixMilli(),
			TypeName:       "exponential_histogram",
		}, {
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key2",
				"value": "value2",
			},
			TimeSeriesValues: map[string][]float64{
				"count": {2.0},
				"sum":   {20.0},
			},
			StartTimestamp: time.Unix(1, 0).UnixMilli(),
			EndTimestamp:   time.Unix(2, 0).UnixMilli(),
			TypeName:       "exponential_histogram",
		},
	},
}

func makeTestcaseFromExpBucket(testcase map[interface{}]interface{}) (pmetric.ExponentialHistogramDataPointSlice, []MetricParquetStruct) {
	dps := pmetric.NewExponentialHistogramDataPointSlice()
	for _, input := range testcase["inputs"].([]map[interface{}]interface{}) {
		dp := dps.AppendEmpty()
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(input["start_timestamp"].(time.Time)))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(input["timestamp"].(time.Time)))
		dp.SetCount(uint64(input["count"].(int)))
		dp.SetSum(float64(input["sum"].(int)))
		attrs := dp.Attributes()
		for k, v := range input["attributes"].(map[interface{}]interface{}) {
			attrs.Insert(k.(string), pcommon.NewValueString(v.(string)))
		}
	}
	return dps, testcase["output"].([]MetricParquetStruct)
}

func TestExpHisto(t *testing.T) {
	raw, expected := makeTestcaseFromExpBucket(ExpHistoTestcase)
	name := "test"
	dps := processExponentialHistogramDataPoints(raw, &name)
	assert.Equal(t, expected, dps)
}

var SummaryTestCase = map[interface{}]interface{}{
	"inputs": []map[interface{}]interface{}{
		{
			"start_timestamp": time.Unix(0, 0),
			"timestamp":       time.Unix(1, 0),
			"quantiles":       []float64{0.5, 0.9, 0.99},
			"values":          []float64{1.0, 2.0, 3.0},
			"attributes": map[interface{}]interface{}{
				"key":   "key1",
				"value": "value1",
			},
		}, {
			"start_timestamp": time.Unix(1, 0),
			"timestamp":       time.Unix(2, 0),
			"quantiles":       []float64{0.5, 0.9, 0.99},
			"values":          []float64{1.0, 2.0, 31.0},
			"attributes": map[interface{}]interface{}{
				"key":   "key2",
				"value": "value2",
			},
		},
	},
	"output": []MetricParquetStruct{
		{
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key1",
				"value": "value1",
			},
			TimeSeriesValues: map[string][]float64{
				"quantiles": {0.5, 0.9, 0.99},
				"values":    {1.0, 2.0, 3.0},
			},
			StartTimestamp: time.Unix(0, 0).UnixMilli(),
			EndTimestamp:   time.Unix(1, 0).UnixMilli(),
			TypeName:       "summary",
		}, {
			MetricName: "test",
			Dimensions: map[string]string{
				"key":   "key2",
				"value": "value2",
			},
			TimeSeriesValues: map[string][]float64{
				"quantiles": {0.5, 0.9, 0.99},
				"values":    {1.0, 2.0, 31.0},
			},
			StartTimestamp: time.Unix(1, 0).UnixMilli(),
			EndTimestamp:   time.Unix(2, 0).UnixMilli(),
			TypeName:       "summary",
		},
	},
}

func makeTestcaseFromSummary(testcase map[interface{}]interface{}) (pmetric.SummaryDataPointSlice, []MetricParquetStruct) {
	dps := pmetric.NewSummaryDataPointSlice()
	for _, input := range testcase["inputs"].([]map[interface{}]interface{}) {
		dp := dps.AppendEmpty()
		dp.SetStartTimestamp(pcommon.NewTimestampFromTime(input["start_timestamp"].(time.Time)))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(input["timestamp"].(time.Time)))
		quantiles := dp.QuantileValues()
		for idx := range input["quantiles"].([]float64) {
			quantile := quantiles.AppendEmpty()
			quantile.SetQuantile(input["quantiles"].([]float64)[idx])
			quantile.SetValue(input["values"].([]float64)[idx])
		}
		attrs := dp.Attributes()
		for k, v := range input["attributes"].(map[interface{}]interface{}) {
			attrs.Insert(k.(string), pcommon.NewValueString(v.(string)))
		}
	}
	return dps, testcase["output"].([]MetricParquetStruct)
}

func TestProcessSummaryDataPoints(t *testing.T) {
	raw, expected := makeTestcaseFromSummary(SummaryTestCase)
	name := "test"
	dps := processSummaryDataPoints(raw, &name)
	assert.Equal(t, expected, dps)
}

func TestFromMetricsToPareut(t *testing.T) {
	results := make([]MetricParquetStruct, 0)
	pms := testdata.GenerateMetricsAllTypesEmptyDataPoint()
	resMetrics := pms.ResourceMetrics()
	for idx := 0; idx < resMetrics.Len(); idx++ {
		scopeMetrics := resMetrics.At(idx).ScopeMetrics()
		for jdx := 0; jdx < scopeMetrics.Len(); jdx++ {
			metrics := scopeMetrics.At(jdx).Metrics()
			results = append(results, FromMetricsToPareut(metrics)...)
		}
	}
	assert.Equal(t, 7, len(results))
}
