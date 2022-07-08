package metadataexporter

const PurviewMetadataSchema = `{
	"entityDefs": [
			{
				"category": "ENTITY",
				"version": 1,
				"name": "metadata",
				"description": "Metadata for opentelemetry data",
				"typeVersion": "1.0.0",
				"attributeDefs": [{
					"name": "destinations",
					"typeName": "array<string>",
					"isOptional": true,
					"cardinality": "LIST",
					"valuesMinCount": 0,
					"valuesMaxCount": 10,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "dimensions",
					"typeName": "map<string,string>",
					"isOptional": true,
					"cardinality": "SINGLE",
					"valuesMinCount": 0,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "sliName",
					"typeName": "string",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "lastPublishedTime",
					"typeName": "date",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "serviceName",
					"typeName": "string",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "sliVersion",
					"typeName": "string",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "serviceGUID",
					"typeName": "string",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}],
				"superTypes": [
					"DataSet"
				],
				"subTypes": [],
				"relationshipAttributeDefs": [

				 ]
			 }
	   ]
	 }`
