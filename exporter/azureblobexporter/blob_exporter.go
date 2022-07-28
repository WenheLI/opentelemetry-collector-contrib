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

package azureblobexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azureblobexporter"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	parquetAzblob "github.com/xitongsys/parquet-go-source/azblob"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type blobExporter struct {
	client             *azblob.ContainerClient
	endpoint           string
	storageAccountName string
	containerName      string
}

func (e *blobExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (e *blobExporter) ConsumeTraces(_ context.Context, td ptrace.Traces) error {
	return errors.New("not implemented")
}

func (e *blobExporter) ConsumeMetrics(_ context.Context, md pmetric.Metrics) error {
	return exportToParquet(e, md)
}

func (e *blobExporter) ConsumeLogs(_ context.Context, ld plog.Logs) error {
	return errors.New("not implemented")
}

func exportToParquet(e *blobExporter, md pmetric.Metrics) error {
	url := fmt.Sprintf("https://%s.%s/%s/%s.parquet", e.storageAccountName, e.endpoint, e.containerName, time.Now().String())
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	config := azblob.ClientOptions{}
	blobWriter, err := parquetAzblob.NewAzBlobFileWriter(context.Background(), url, cred, config)
	if err != nil {
		return err
	}

	parquetWriter, err := writer.NewJSONWriter(JSONSCHEMA, blobWriter, 4)
	if err != nil {
		return err
	}

	parquetWriter.RowGroupSize = 128 * 1024 * 1024 //128M
	parquetWriter.PageSize = 8 * 1024              //8K
	parquetWriter.CompressionType = parquet.CompressionCodec_SNAPPY

	parquetData := make([]MetricParquetStruct, 0)
	metrics := md.ResourceMetrics()
	for i := 0; i < metrics.Len(); i++ {
		scopedMetrics := metrics.At(i)
		for j := 0; j < scopedMetrics.ScopeMetrics().Len(); j++ {
			metric := scopedMetrics.ScopeMetrics().At(j).Metrics()
			parquetData = append(parquetData, FromMetricsToPareut(metric)...)
		}
	}

	for _, metric := range parquetData {
		content, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		err = parquetWriter.Write(content)
		if err != nil {
			return err
		}
	}

	err = parquetWriter.WriteStop()
	if err != nil {
		return err
	}
	err = blobWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (e *blobExporter) Start(context.Context, component.Host) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s.%s", e.storageAccountName, e.endpoint)
	client, err := azblob.NewServiceClient(url, cred, nil)
	if err != nil {
		return err
	}
	listContainersOptions := azblob.ListContainersOptions{
		Include: azblob.ListContainersDetail{
			Metadata: true, // Include Metadata
			Deleted:  true, // Include deleted containers in the result as well
		},
	}
	pager := client.ListContainers(&listContainersOptions)
	containers := make(map[string]bool)
	for pager.NextPage(context.TODO()) {
		resp := pager.PageResponse()

		for _, container := range resp.ContainerItems {
			containers[*container.Name] = true
		}
	}
	containerClient, err := client.NewContainerClient(e.containerName)
	if err != nil {
		return err
	}

	if !containers[e.containerName] {
		ctx := context.Background()
		_, err = containerClient.Create(ctx, nil)
		if err != nil {
			return err
		}
	}

	e.client = containerClient

	return nil
}

func (e *blobExporter) Shutdown(context.Context) error {
	e.client = nil
	return nil
}
