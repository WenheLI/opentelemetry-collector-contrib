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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
)

func TestConsumeTraces(t *testing.T) {
	exporter := &metadataExporter{}

	require.NotNil(t, exporter)
	td := testdata.GenerateTracesTwoSpansSameResource()
	err := exporter.ConsumeTraces(context.Background(), td)
	require.EqualError(t, err, "not implemented")
}

func TestConsumeLogs(t *testing.T) {
	exporter := &metadataExporter{}
	require.NotNil(t, exporter)
	td := testdata.GenerateLogsTwoLogRecordsSameResource()
	err := exporter.ConsumeLogs(context.Background(), td)
	require.EqualError(t, err, "not implemented")
}

// mock a client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) authentication() error {
	return nil
}

func (m *MockClient) CheckMetadataType() (bool, error) {
	return true, nil
}

func (m *MockClient) CreateMetadataType() (bool, error) {
	return true, nil
}

func (m *MockClient) CreateMetadataEntity(entities PurviewEntityBulkType) (bool, error) {
	return true, nil
}

func TestConsumeMetrics(t *testing.T) {
	mockClient := new(MockClient)
	exporter := &metadataExporter{
		destinations: []string{"dummyDestination"},
		endpoint:     "dummyEndpoint",
		accountName:  "dummyAccountName",
		client:       mockClient,
	}
	require.NotNil(t, exporter)
	td := testdata.GenerateMetricsTwoMetrics()
	assert.NoError(t, exporter.Start(context.Background(), componenttest.NewNopHost()))
	assert.NoError(t, exporter.ConsumeMetrics(context.Background(), td))
	assert.NoError(t, exporter.Shutdown(context.Background()))
}
