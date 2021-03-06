/*
Copyright 2018 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewClusterResource(t *testing.T) {
	for _, c := range []struct {
		desc     string
		resource *PipelineResource
		want     *ClusterResource
	}{{
		desc: "basic cluster resource",
		resource: &PipelineResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-resource",
				Namespace: "foo",
			},
			Spec: PipelineResourceSpec{
				Type: PipelineResourceTypeCluster,
				Params: []Param{{
					Name:  "url",
					Value: "http://10.10.10.10",
				}, {
					Name:  "cadata",
					Value: "bXktY2x1c3Rlci1jZXJ0Cg",
				}, {
					Name:  "token",
					Value: "my-token",
				},
				},
			},
		},
		want: &ClusterResource{
			Name:   "test-cluster-resource",
			Type:   PipelineResourceTypeCluster,
			URL:    "http://10.10.10.10",
			CAData: []byte("my-cluster-cert"),
			Token:  "my-token",
		},
	}, {
		desc: "resource with password instead of token",
		resource: &PipelineResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-resource",
				Namespace: "foo",
			},
			Spec: PipelineResourceSpec{
				Type: PipelineResourceTypeCluster,
				Params: []Param{{
					Name:  "url",
					Value: "http://10.10.10.10",
				}, {
					Name:  "cadata",
					Value: "bXktY2x1c3Rlci1jZXJ0Cg",
				}, {
					Name:  "username",
					Value: "user",
				}, {
					Name:  "password",
					Value: "pass",
				},
				},
			},
		},
		want: &ClusterResource{
			Name:     "test-cluster-resource",
			Type:     PipelineResourceTypeCluster,
			URL:      "http://10.10.10.10",
			CAData:   []byte("my-cluster-cert"),
			Username: "user",
			Password: "pass",
		},
	}, {
		desc: "set insecure flag to true when there is no cert",
		resource: &PipelineResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-resource",
				Namespace: "foo",
			},
			Spec: PipelineResourceSpec{
				Type: PipelineResourceTypeCluster,
				Params: []Param{{
					Name:  "url",
					Value: "http://10.10.10.10",
				}, {
					Name:  "token",
					Value: "my-token",
				},
				},
			},
		},
		want: &ClusterResource{
			Name:     "test-cluster-resource",
			Type:     PipelineResourceTypeCluster,
			URL:      "http://10.10.10.10",
			Token:    "my-token",
			Insecure: true,
		},
	}, {
		desc: "basic resource with secrets",
		resource: &PipelineResource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster-resource",
				Namespace: "foo",
			},
			Spec: PipelineResourceSpec{
				Type: PipelineResourceTypeCluster,
				Params: []Param{{
					Name:  "url",
					Value: "http://10.10.10.10",
				}},
				SecretParams: []SecretParam{{
					FieldName:  "cadata",
					SecretKey:  "cadatakey",
					SecretName: "secret1",
				}, {
					FieldName:  "token",
					SecretKey:  "tokenkey",
					SecretName: "secret1",
				}},
			},
		},
		want: &ClusterResource{
			Name: "test-cluster-resource",
			Type: PipelineResourceTypeCluster,
			URL:  "http://10.10.10.10",
			Secrets: []SecretParam{{
				FieldName:  "cadata",
				SecretKey:  "cadatakey",
				SecretName: "secret1",
			}, {
				FieldName:  "token",
				SecretKey:  "tokenkey",
				SecretName: "secret1",
			}},
		},
	}} {
		t.Run(c.desc, func(t *testing.T) {
			got, err := NewClusterResource(c.resource)
			if err != nil {
				t.Errorf("Test: %q; TestNewClusterResource() error = %v", c.desc, err)
			}
			if d := cmp.Diff(got, c.want); d != "" {
				t.Errorf("Diff:\n%s", d)
			}
		})
	}
}
