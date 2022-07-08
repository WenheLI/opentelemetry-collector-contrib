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

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestExtractEmptyResource(t *testing.T) {
	resource := pcommon.NewResource()
	result := extractResource(resource)
	assert.Empty(t, result)
}

type MockResource = map[string]interface{}
type ExpectResult = map[string]string

type TestCase struct {
	mockResource MockResource
	expectResult ExpectResult
}

var testCases = []TestCase{
	{MockResource{"resource-attr": "resource-attr-val-1"}, ExpectResult{"resource-attr": "resource-attr-val-1"}},
	{MockResource{"resource-attr": ""}, ExpectResult{"resource-attr": ""}},
}

func TestExtractResouce(t *testing.T) {
	for _, testCase := range testCases {
		resource := pcommon.NewResource()
		for k, v := range testCase.mockResource {
			resource.Attributes().UpsertString(k, v.(string))
		}
		result := extractResource(resource)
		assert.Equal(t, testCase.expectResult, result)
	}
}
