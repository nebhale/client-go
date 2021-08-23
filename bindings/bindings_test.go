/*
 * Copyright 2021 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bindings_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/nebhale/client-go/bindings"
)

func TestCached(t *testing.T) {
	b := bindings.Cached([]bindings.Binding{
		bindings.MapBinding{Name: "test-name-1"},
		bindings.MapBinding{Name: "test-name-2"},
	})

	for _, c := range b {
		if _, ok := c.(*bindings.CacheBinding); !ok {
			t.Errorf("did not wrap with CacheBinding")
		}
	}
}

func TestFrom_Invalid(t *testing.T) {
	if !reflect.DeepEqual(bindings.From("invalid"), []bindings.Binding{}) {
		t.Errorf("did not create an empty Bindings")
	}
}

func TestFrom_NotDirectory(t *testing.T) {
	if !reflect.DeepEqual(bindings.From("testdata/additional-file"), []bindings.Binding{}) {
		t.Errorf("did not create an empty Bindings")
	}
}

func TestFrom_Exists(t *testing.T) {
	if len(bindings.From("testdata")) != 3 {
		t.Errorf("did not create proper number of bindings")
	}
}

func TestFromServiceBindingRoot_Unset(t *testing.T) {
	if !reflect.DeepEqual(bindings.FromServiceBindingRoot(), []bindings.Binding{}) {
		t.Errorf("did not create an empty Bindings")
	}
}

func TestFromServiceBindingRoot_Set(t *testing.T) {
	old, ok := os.LookupEnv("SERVICE_BINDING_ROOT")
	if err := os.Setenv("SERVICE_BINDING_ROOT", "testdata"); err != nil {
		t.Errorf("unable to set SERVICE_BINDING_ROOT")
	}
	defer func() {
		if ok {
			if err := os.Setenv("SERVICE_BINDING_ROOT", old); err != nil {
				t.Errorf("unable to set SERVICE_BINDING_ROOT")
			}
		} else {
			if err := os.Unsetenv("SERVICE_BINDING_ROOT"); err != nil {
				t.Errorf("unable to unset SERVICE_BINDING_ROOT")
			}
		}
	}()

	if len(bindings.FromServiceBindingRoot()) != 3 {
		t.Errorf("did not create proper number of bindings")
	}
}

func TestFind_Invalid(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{Name: "test-name-1"},
	}

	if _, ok := bindings.Find(b, "test-name-2"); ok {
		t.Errorf("does not identify invalid name")
	}
}

func TestFind_Valid(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{Name: "test-name-1"},
		bindings.MapBinding{Name: "test-name-2"},
	}

	if c, ok := bindings.Find(b, "test-name-1"); !ok {
		t.Errorf("does not identify valid name")
	} else if c.GetName() != "test-name-1" {
		t.Errorf("does not return valid binding")
	}
}

func TestFilter_Empty(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type": []byte("test-type-2"),
			},
		},
	}

	if len(bindings.Filter(b, "test-type-3")) != 0 {
		t.Errorf("incorrect number of matches")
	}
}

func TestFilter_Single(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type": []byte("test-type-2"),
			},
		},
	}

	if len(bindings.Filter(b, "test-type-1")) != 1 {
		t.Errorf("incorrect number of matches")
	}
}

func TestFilter_Multiple(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type": []byte("test-type-1"),
			},
		},
	}

	if len(bindings.Filter(b, "test-type-1")) != 2 {
		t.Errorf("incorrect number of matches")
	}
}

func TestFilterWithProvider_Empty(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type":     []byte("test-type-2"),
				"provider": []byte("test-provider-2"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-3",
			Content: map[string][]byte{
				"type": []byte("test-type-3"),
			},
		},
	}

	if len(bindings.FilterWithProvider(b, "test-type-1", "test-provider-2")) != 0 {
		t.Errorf("incorrect number of matches")
	}
}

func TestFilterWithProvider_Single(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-2"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-3",
			Content: map[string][]byte{
				"type": []byte("test-type-3"),
			},
		},
	}

	if len(bindings.FilterWithProvider(b, "test-type-1", "test-provider-2")) != 1 {
		t.Errorf("incorrect number of matches")
	}
}

func TestFilterWithProvider_Multiple(t *testing.T) {
	b := []bindings.Binding{
		bindings.MapBinding{
			Name: "test-name-1",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-2",
			Content: map[string][]byte{
				"type":     []byte("test-type-1"),
				"provider": []byte("test-provider-1"),
			},
		},
		bindings.MapBinding{
			Name: "test-name-3",
			Content: map[string][]byte{
				"type": []byte("test-type-3"),
			},
		},
	}

	if len(bindings.FilterWithProvider(b, "test-type-1", "test-provider-1")) != 2 {
		t.Errorf("incorrect number of matches")
	}
}
