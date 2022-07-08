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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
)

func TestConsumeTraces(t *testing.T) {
	f, err := os.CreateTemp("", "*.json")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	exporter := &metadataExporter{
		path: f.Name(),
	}

	require.NotNil(t, exporter)
	td := testdata.GenerateTracesTwoSpansSameResource()
	err = exporter.ConsumeTraces(context.Background(), td)
	require.EqualError(t, err, "not implemented")
}

func TestConsumeLogs(t *testing.T) {
	f, err := os.CreateTemp("", "*.json")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	exporter := &metadataExporter{
		path: f.Name(),
	}
	require.NotNil(t, exporter)
	td := testdata.GenerateLogsTwoLogRecordsSameResource()
	err = exporter.ConsumeLogs(context.Background(), td)
	require.EqualError(t, err, "not implemented")
}

func TestConsumeMetrics(t *testing.T) {
	f, err := os.CreateTemp("", "*.json")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	exporter := &metadataExporter{
		path:         f.Name(),
		destinations: []string{"dummyDestination"},
	}
	require.NotNil(t, exporter)
	td := testdata.GenerateMetricsTwoMetrics()
	assert.NoError(t, exporter.Start(context.Background(), componenttest.NewNopHost()))
	assert.NoError(t, exporter.ConsumeMetrics(context.Background(), td))
	assert.NoError(t, exporter.Shutdown(context.Background()))
	buf, err := ioutil.ReadFile(exporter.path)
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	var ret []MetricMetadata
	err = json.Unmarshal(buf, &ret)
	require.NoError(t, err)
	require.Equal(t, 1, len(ret))
	assert.Equal(t, "dummyDestination", ret[0].Destinations[0])
}
