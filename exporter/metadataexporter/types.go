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

package metadataexporter

import "fmt"

type MetricMetadataPoint struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	LastPublishtime int64             `json:"lastPublishtime"`
	Dimensions      map[string]string `json:"dimensions"` // key dimension name, value dimension data type
}

type MetricMetadata struct {
	MetricMetadataPoints []MetricMetadataPoint        `json:"metadata"`
	Resources            map[string]string            `json:"resources"`
	Destinations         map[string]map[string]string `json:"destinations"`
}

type PurviewAttributes struct {
	QualifiedName     string                       `json:"qualifiedName"`
	Name              string                       `json:"name"`
	Description       string                       `json:"description"`
	PrincipalId       int                          `json:"principalId"`
	LastPublishedTime int64                        `json:"lastPublishedTime"`
	Dimensions        map[string]string            `json:"dimensions"`
	Destinations      map[string]map[string]string `json:"destinations"`
	ObjectType        string                       `json:"objectType"`
	ServiceName       string                       `json:"serviceName"`
	ServiceGUID       string                       `json:"serviceGUID"`
	SLIName           string                       `json:"sliName"`
	SLIVersion        string                       `json:"sliVersion"`
}

type PurviewMetadataEntity struct {
	Meanings   []string          `json:"meanings"`
	Status     string            `json:"status"`
	Version    int               `json:"version"`
	TypeName   string            `json:"typeName"`
	Attributes PurviewAttributes `json:"attributes"`
}

type PurviewEntityBulkType struct {
	ReferredEntities map[string]interface{}  `json:"referredEntities"`
	Entities         []PurviewMetadataEntity `json:"entities"`
}

func NewPurviewEntity(point MetricMetadataPoint, resource map[string]string, destinations map[string]map[string]string) PurviewMetadataEntity {
	qualifiedName := fmt.Sprintf("%s-%s", resource["service.name"], point.Name)
	purviewAttributes := PurviewAttributes{
		QualifiedName:     qualifiedName,
		Name:              point.Name,
		Description:       point.Description,
		PrincipalId:       0,
		LastPublishedTime: point.LastPublishtime,
		Dimensions:        point.Dimensions,
		Destinations:      destinations,
		ObjectType:        "",
		SLIName:           point.Dimensions["sliName"],
		SLIVersion:        point.Dimensions["sliVersion"],
		ServiceName:       resource["service.name"],
		ServiceGUID:       resource["service.instance.id"],
	}

	// reset to default value if it is empty
	if purviewAttributes.SLIName == "" {
		purviewAttributes.SLIName = "default"
	}
	if purviewAttributes.SLIVersion == "" {
		purviewAttributes.SLIVersion = "default"
	}

	purviewMetadataEntity := PurviewMetadataEntity{
		Meanings:   make([]string, 0),
		Status:     "ACTIVE",
		Version:    0,
		TypeName:   "metadata",
		Attributes: purviewAttributes,
	}

	return purviewMetadataEntity
}
