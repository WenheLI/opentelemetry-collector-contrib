package azureblobexporter

import (
	"encoding/json"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

type MetricParquetStruct struct {
	StartTimestamp   int64                `json:"Start Timestamp"`
	EndTimestamp     int64                `json:"End Timestamp"`
	MetricName       string               `json:"Metric Name"`
	TimeSeriesValues map[string][]float64 `json:"Time Series Values"`
	Dimensions       map[string]string    `json:"Dimensions"`
	TypeName         string               `json:"Type Name"`
}

var SCHEMA = arrow.NewSchema(
	[]arrow.Field{
		{Name: "Start Timestamp", Type: arrow.FixedWidthTypes.Timestamp_ms},
		{Name: "End Timestamp", Type: arrow.FixedWidthTypes.Timestamp_ms},
		{Name: "Metric Name", Type: arrow.BinaryTypes.String},
		{Name: "Dimensions", Type: arrow.MapOf(
			arrow.BinaryTypes.String,
			arrow.BinaryTypes.String,
		)},
		{Name: "Type Name", Type: arrow.BinaryTypes.String},
		{Name: "Time Series Values", Type: arrow.MapOf(
			arrow.BinaryTypes.String,
			arrow.ListOf(arrow.PrimitiveTypes.Float64),
		)},
	},
	nil,
)

func BuildRecordFrom(dataset []MetricParquetStruct) array.Record {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	b := array.NewRecordBuilder(mem, SCHEMA)
	defer b.Release()
	fields := SCHEMA.Fields()
	for _, data := range dataset {
		for idx := 0; idx < len(fields); idx++ {
			switch fields[idx].Name {
			case "Start Timestamp":
				b.Field(idx).(*array.Int64Builder).Append(data.StartTimestamp)
			case "End Timestamp":
				b.Field(idx).(*array.Int64Builder).Append(data.EndTimestamp)
			case "Metric Name":
				b.Field(idx).(*array.StringBuilder).Append(data.MetricName)
			case "Dimensions":
				mb := b.Field(idx).(*array.MapBuilder)
				kb := mb.KeyBuilder().(*array.StringBuilder)
				vb := mb.ItemBuilder().(*array.StringBuilder)
				mb.Append(true)
				for k, v := range data.Dimensions {
					kb.Append(k)
					vb.Append(v)
				}
			case "Type Name":
				b.Field(idx).(*array.StringBuilder).Append(data.TypeName)
			case "Time Series Values":
				mb := b.Field(idx).(*array.MapBuilder)
				kb := mb.KeyBuilder().(*array.StringBuilder)
				vb := mb.ItemBuilder().(*array.ListBuilder)
				mb.Append(true)
				for k, v := range data.TimeSeriesValues {
					kb.Append(k)
					vvb := vb.ValueBuilder().(*array.Float64Builder)
					vb.Append(true)
					for _, vv := range v {
						vvb.Append(vv)
					}
				}
			}
		}
	}
	return b.NewRecord()
}

const JSONSCHEMA = `
{
	"Tag":"name=parquet-go-root",
	"Fields": [
		{"Tag": "name=Start Timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"},
		{"Tag": "name=End Timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"},
		{"Tag": "name=Metric Name, type=BYTE_ARRAY, convertedtype=UTF8"},
		{"Tag": "name=Dimensions, type=MAP", "Fields": [
			{"Tag": "name=key, type=BYTE_ARRAY, convertedtype=UTF8"},
			{"Tag": "name=value, type=BYTE_ARRAY, convertedtype=UTF8"}	
		]},
		{"Tag": "name=Type Name, type=BYTE_ARRAY, convertedtype=UTF8"},
		{"Tag": "name=Time Series Values, type=MAP", "Fields": [
			{"Tag": "name=key, type=BYTE_ARRAY, convertedtype=UTF8"},
			{"Tag": "name=value, type=LIST", "Fields": [
				{"Tag": "name=element, type=DOUBLE"}
			]}
		]}
	]
}
`

func BuildJSONFrom(dataset MetricParquetStruct) string {
	bytes, _ := json.Marshal(dataset)
	return string(bytes)
}
