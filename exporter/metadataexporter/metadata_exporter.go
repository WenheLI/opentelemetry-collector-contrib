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

/*
 * Main entrance for consuming incoming OTEL data flow and converting it to metadata format
 */

package metadataexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/metadataexporter"

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type metadataExporter struct {
	path  string
	file  io.WriteCloser
	mutex sync.Mutex
}

func (e *metadataExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *metadataExporter) ConsumeTraces(_ context.Context, td ptrace.Traces) error {
	panic("not implemented")
}

func (e *metadataExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	var MetricMetadataPoints []MetricMetadataPoint
	var metricMetadataList = make([]MetricMetadata, 0)

	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resouceMetric := resourceMetrics.At(i)
		// extract resources for the metric
		resources := extractResource(resouceMetric.Resource())

		// extract metrics metadata
		MetricMetadataPoints = extractMetric(resouceMetric.ScopeMetrics())
		metricMetadata := MetricMetadata{
			Resources:            resources,
			MetricMetadataPoints: MetricMetadataPoints,
		}

		metricMetadataList = append(metricMetadataList, metricMetadata)
	}

	content, _ := json.Marshal(metricMetadataList)
	return exportMessageAsLine(e, content)
}

func (e *metadataExporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	panic("not implemented")
}

func exportMessageAsLine(e *metadataExporter, buf []byte) error {
	// Ensure only one write operation happens at a time.
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if _, err := e.file.Write(buf); err != nil {
		return err
	}
	if _, err := io.WriteString(e.file, "\n"); err != nil {
		return err
	}
	return nil
}

func (e *metadataExporter) Start(context.Context, component.Host) error {
	var err error
	e.file, err = os.OpenFile(e.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	return err
}

// Shutdown stops the exporter and is invoked during shutdown.
func (e *metadataExporter) Shutdown(context.Context) error {
	return e.file.Close()
}
