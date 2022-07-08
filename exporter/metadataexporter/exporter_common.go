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
 * Utility functions to extract metadata from OTEL metrics(For resource data)
 */

package metadataexporter

import "go.opentelemetry.io/collector/pdata/pcommon"

func extractResource(resource pcommon.Resource) map[string]string {
	var data = make(map[string]string)
	attrs := resource.Attributes()
	attrs.Range(func(k string, v pcommon.Value) bool {
		data[k] = v.StringVal()
		return true
	})
	return data
}
