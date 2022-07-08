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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPurviewTransport(t *testing.T) {
	pt := &PurviewTransport{
		token: "token",
	}
	assert.Equal(t, "token", pt.token)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	pt.RoundTrip(req)
	assert.Equal(t, "Bearer token", req.Header.Get("Authorization"))
}

func TestNewPurviewClient(t *testing.T) {
	pc, err := NewPurviewClient("endpoint", "account")
	assert.NoError(t, err)
	assert.Equal(t, "endpoint", pc.endpoint)
	assert.Equal(t, "account", pc.accountName)
}
