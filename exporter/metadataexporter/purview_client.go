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

// Client for Purview API
package metadataexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const (
	PurviewMetadataTypeName = "Metadata"
	PurviewTypedefNameAPI   = "/api/atlas/v2/types/typedef/name/"
	PurviewTypedefAPI       = "/api/atlas/v2/types/typedefs/"
	PurviewEntityAPI        = "/api/atlas/v2/entity/bulk/"
)

type IPurviewClient interface {
	CheckMetadataType() (bool, error)
	CreateMetadataType() (bool, error)
	CreateMetadataEntity(entities PurviewEntityBulkType) (bool, error)
	authentication() error
}

// for testing
type IAZIdentityCredential interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type PurviewClient struct {
	identity    IAZIdentityCredential
	httpClient  *http.Client
	endpoint    string
	accountName string
}

type PurviewTransport struct {
	token string
}

func (t *PurviewTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", t.token))
	return http.DefaultTransport.RoundTrip(req)
}

func NewPurviewClient(endpoint string, accountName string) (*PurviewClient, error) {
	httpClient := &http.Client{}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	httpClient.Transport = &PurviewTransport{}
	return &PurviewClient{
		identity:    cred,
		httpClient:  httpClient,
		endpoint:    endpoint,
		accountName: accountName,
	}, nil
}

func (client *PurviewClient) authentication() error {
	opts := policy.TokenRequestOptions{
		Scopes: []string{"https://purview.azure.net/.default"},
	}
	token, err := client.identity.GetToken(context.Background(), opts)
	if err != nil {
		return err
	}
	client.httpClient.Transport.(*PurviewTransport).token = token.Token
	return nil
}

func (client *PurviewClient) CheckMetadataType() (bool, error) {
	err := client.authentication()
	if err != nil {
		return false, err
	}
	req, err := client.httpClient.Get(fmt.Sprint(client.endpoint, PurviewTypedefNameAPI, PurviewMetadataTypeName))
	if err != nil {
		return false, err
	}
	defer req.Body.Close()
	if req.StatusCode/100 == 2 {
		return true, nil
	}
	return false, nil
}

func (client *PurviewClient) CreateMetadataType() (bool, error) {
	err := client.authentication()
	if err != nil {
		return false, err
	}
	result, err := client.CheckMetadataType()
	if err != nil {
		return false, err
	}
	if result {
		return true, nil
	}
	req, err := client.httpClient.Post(fmt.Sprint(client.endpoint, PurviewTypedefAPI),
		"application/json", bytes.NewBuffer([]byte(PurviewMetadataSchema)))
	if err != nil {
		return false, err
	}
	defer req.Body.Close()
	if req.StatusCode/100 == 2 {
		return true, nil
	}
	return false, nil
}

func (client *PurviewClient) CreateMetadataEntity(entities PurviewEntityBulkType) (bool, error) {
	err := client.authentication()
	if err != nil {
		return false, err
	}

	buf, err := json.Marshal(entities)
	if err != nil {
		return false, err
	}
	req, err := client.httpClient.Post(fmt.Sprint(client.endpoint, PurviewEntityAPI),
		"application/json", bytes.NewBuffer(buf))
	if err != nil {
		return false, err
	}
	if req.StatusCode/100 != 2 {
		return false, nil
	}

	return true, nil
}
