package azureblobexporter

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func processHistogramDataPoints(dps pmetric.HistogramDataPointSlice, name *string) []MetricParquetStruct {
	ret := make([]MetricParquetStruct, 0)
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dimensions := make(map[string]string)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			dimensions[k] = v.AsString()
			return true
		})
		bounds := dp.ExplicitBounds()
		counts := dp.BucketCounts()

		boundsList := make([]float64, 0)
		countsList := make([]float64, 0)

		for i := 0; i < bounds.Len(); i++ {
			bound := bounds.At(i)
			boundsList = append(boundsList, bound)
		}

		for i := 0; i < counts.Len(); i++ {
			count := counts.At(i)
			// overflow
			countsList = append(countsList, float64(count))
		}

		ret = append(ret, MetricParquetStruct{
			StartTimestamp: dp.StartTimestamp().AsTime().UnixMilli(),
			EndTimestamp:   dp.Timestamp().AsTime().UnixMilli(),
			MetricName:     *name,
			Dimensions:     dimensions,
			TypeName:       "histogram",
			TimeSeriesValues: map[string][]float64{
				"bounds": boundsList,
				"counts": countsList,
			},
		})
	}
	return ret
}

func processNumberDataPoints(dps pmetric.NumberDataPointSlice, name *string) []MetricParquetStruct {
	ret := make([]MetricParquetStruct, 0)
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dimensions := make(map[string]string)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			dimensions[k] = v.AsString()
			return true
		})
		var value float64
		if dp.ValueType() == pmetric.NumberDataPointValueTypeInt {
			value = float64(dp.IntVal())
		} else {
			value = float64(dp.DoubleVal())
		}
		ret = append(ret, MetricParquetStruct{
			StartTimestamp: dp.StartTimestamp().AsTime().UnixMilli(),
			EndTimestamp:   dp.Timestamp().AsTime().UnixMilli(),
			MetricName:     *name,
			Dimensions:     dimensions,
			TypeName:       "number",
			TimeSeriesValues: map[string][]float64{
				"value": {value},
			},
		})
	}
	return ret
}

func processExponentialHistogramDataPoints(dps pmetric.ExponentialHistogramDataPointSlice, name *string) []MetricParquetStruct {
	ret := make([]MetricParquetStruct, 0)
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dimensions := make(map[string]string)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			dimensions[k] = v.AsString()
			return true
		})
		count := float64(dp.Count())
		sum := float64(dp.Sum())
		ret = append(ret, MetricParquetStruct{
			StartTimestamp: dp.StartTimestamp().AsTime().UnixMilli(),
			EndTimestamp:   dp.Timestamp().AsTime().UnixMilli(),
			MetricName:     *name,
			Dimensions:     dimensions,
			TypeName:       "exponential_histogram",
			TimeSeriesValues: map[string][]float64{
				"count": {count},
				"sum":   {sum},
			},
		})
	}
	return ret
}

func processSummaryDataPoints(dps pmetric.SummaryDataPointSlice, name *string) []MetricParquetStruct {
	ret := make([]MetricParquetStruct, 0)
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		dimensions := make(map[string]string)
		dp.Attributes().Range(func(k string, v pcommon.Value) bool {
			dimensions[k] = v.AsString()
			return true
		})

		quantiles := make([]float64, 0)
		values := make([]float64, 0)
		for i := 0; i < dp.QuantileValues().Len(); i++ {
			quantile := dp.QuantileValues().At(i).Quantile()
			value := dp.QuantileValues().At(i).Value()
			quantiles = append(quantiles, quantile)
			values = append(values, value)
		}
		ret = append(ret, MetricParquetStruct{
			StartTimestamp: dp.StartTimestamp().AsTime().UnixMilli(),
			EndTimestamp:   dp.Timestamp().AsTime().UnixMilli(),
			MetricName:     *name,
			Dimensions:     dimensions,
			TypeName:       "summary",
			TimeSeriesValues: map[string][]float64{
				"quantiles": quantiles,
				"values":    values,
			},
		})
	}
	return ret
}

func FromMetricsToPareut(pm pmetric.MetricSlice) []MetricParquetStruct {
	ret := make([]MetricParquetStruct, 0)
	for i := 0; i < pm.Len(); i++ {
		metric := pm.At(i)
		metricName := metric.Name()
		switch metric.DataType() {
		case pmetric.MetricDataTypeGauge:
			ret = append(ret, processNumberDataPoints(metric.Gauge().DataPoints(), &metricName)...)
		case pmetric.MetricDataTypeSum:
			ret = append(ret, processNumberDataPoints(metric.Sum().DataPoints(), &metricName)...)
		case pmetric.MetricDataTypeHistogram:
			ret = append(ret, processHistogramDataPoints(metric.Histogram().DataPoints(), &metricName)...)
		case pmetric.MetricDataTypeExponentialHistogram:
			ret = append(ret, processExponentialHistogramDataPoints(metric.ExponentialHistogram().DataPoints(), &metricName)...)
		case pmetric.MetricDataTypeSummary:
			ret = append(ret, processSummaryDataPoints(metric.Summary().DataPoints(), &metricName)...)
		}
	}

	return ret
}
