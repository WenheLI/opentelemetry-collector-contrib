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

const PurviewMetadataSchema = `{
	"entityDefs": [
			{
				"category": "ENTITY",
				"version": 2,
				"name": "metadata",
				"description": "Metadata for opentelemetry data",
				"typeVersion": "2.0.0",
				"attributeDefs": [{
					"name": "destinations",
					"typeName": "map<string,map<string,string>>",
					"isOptional": true,
					"cardinality": "SINGLE",
					"valuesMinCount": 0,
					"valuesMaxCount": 1,
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
				}, {
					"name": "isSLI",
					"typeName": "boolean",
					"isOptional": false,
					"cardinality": "SINGLE",
					"valuesMinCount": 1,
					"valuesMaxCount": 1,
					"isUnique": false,
					"isIndexable": true,
					"includeInNotification": false
				}, {
					"name": "sliDetails",
					"typeName": "map<string,string>",
					"isOptional": true,
					"cardinality": "SINGLE",
					"valuesMinCount": 0,
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
