package azureblobexporter

type MetricParquetStruct struct {
	StartTimestamp   int64                `json:"Start Timestamp"`
	EndTimestamp     int64                `json:"End Timestamp"`
	MetricName       string               `json:"Metric Name"`
	TimeSeriesValues map[string][]float64 `json:"Time Series Values"`
	Dimensions       map[string]string    `json:"Dimensions"`
	TypeName         string               `json:"Type Name"`
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
