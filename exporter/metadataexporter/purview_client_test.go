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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

type MockIdentity struct{}

func (m *MockIdentity) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "token"}, nil
}

func TestAuthentication(t *testing.T) {
	pc, err := NewPurviewClient("endpoint", "account")
	pc.identity = &MockIdentity{}
	assert.NoError(t, err)
	err = pc.authentication()
	assert.NoError(t, err)
	assert.Equal(t, pc.httpClient.Transport.(*PurviewTransport).token, "token")
}

func TestCheckMetadataType(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/api/atlas/v2/types/typedef/name/Metadata", req.URL.String())
		res.WriteHeader(http.StatusOK)
	}))

	defer testServer.Close()
	pc, err := NewPurviewClient(testServer.URL, "account")
	pc.identity = &MockIdentity{}
	assert.NoError(t, err)
	res, err := pc.CheckMetadataType()
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestCheckMetadataTypeBadReq(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/api/atlas/v2/types/typedef/name/Metadata", req.URL.String())
		res.WriteHeader(http.StatusBadGateway)
	}))

	defer testServer.Close()
	pc, err := NewPurviewClient(testServer.URL, "account")
	pc.identity = &MockIdentity{}
	assert.NoError(t, err)
	res, err := pc.CheckMetadataType()
	assert.NoError(t, err)
	assert.Equal(t, false, res)
}

func TestCreateMetadataType(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		switch req.URL.String() {
		case "/api/atlas/v2/types/typedef/name/":
			assert.Equal(t, "GET", req.Method)
			res.WriteHeader(http.StatusOK)
		case "/api/atlas/v2/types/typedef/":
			assert.Equal(t, "POST", req.Method)
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, PurviewMetadataSchema, string(body))
			res.WriteHeader(http.StatusOK)
		}
	}))
	defer testServer.Close()

	pc, err := NewPurviewClient(testServer.URL, "account")
	pc.identity = &MockIdentity{}
	assert.NoError(t, err)
	res, err := pc.CreateMetadataType()
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestCreateMetadataEntity(t *testing.T) {
	mockEntities := PurviewEntityBulkType{
		Entities: []PurviewMetadataEntity{{
			Meanings: []string{"meaning1", "meaning2"},
			Status:   "status",
		}},
		ReferredEntities: map[string]interface{}{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "/api/atlas/v2/entity/bulk/", req.URL.String())
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		buf, err := json.Marshal(mockEntities)
		assert.NoError(t, err)
		assert.Equal(t, string(buf), string(body))
	}))
	defer testServer.Close()

	pc, err := NewPurviewClient(testServer.URL, "account")
	pc.identity = &MockIdentity{}
	assert.NoError(t, err)
	res, err := pc.CreateMetadataEntity(mockEntities)
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
