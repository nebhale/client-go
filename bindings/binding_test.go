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
	"bytes"
	"path/filepath"
	"testing"

	"github.com/nebhale/client-go/bindings"
)

func TestGet(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"test-secret-key": []byte("test-secret-value\n"),
		},
	}

	if v, ok := bindings.Get(b, "test-secret-key"); !ok {
		t.Errorf("does not identify valid key")
	} else if v != "test-secret-value" {
		t.Errorf("returned the wrong value")
	}
}

func TestGetProvider_InvalidKey(t *testing.T) {
	b := bindings.MapBinding{
		Name:    "test-name",
		Content: map[string][]byte{},
	}

	if _, ok := bindings.GetProvider(b); ok {
		t.Errorf("does not identify invalid key")
	}
}

func TestGetProvider_ValidKey(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"provider": []byte("test-provider-1"),
		},
	}

	if v, ok := bindings.GetProvider(b); !ok {
		t.Errorf("does not identify valid key")
	} else if v != "test-provider-1" {
		t.Errorf("returned the wrong value")
	}
}

func TestGetType_InvalidBinding(t *testing.T) {
	b := bindings.MapBinding{
		Name:    "test-name",
		Content: map[string][]byte{},
	}

	if _, err := bindings.GetType(b); err == nil {
		t.Errorf("does not identify invalid binding")
	}

}

func TestGetType_ValidBinding(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"type": []byte("test-type-1"),
		},
	}

	if v, err := bindings.GetType(b); err != nil {
		t.Errorf("does not identify valid binding")
	} else if v != "test-type-1" {
		t.Errorf("returned the wrong value")
	}
}

func TestCacheBinding_Uncached(t *testing.T) {
	s := &stubBinding{}
	b := bindings.CacheBinding{Delegate: s}

	if v, ok := b.GetAsBytes("test-key"); !ok || v == nil {
		t.Errorf("did not retrieve value as bytes")
	}
	if s.getAsBytesCount != 1 {
		t.Errorf("did not call delegate enough")
	}
}

func TestCacheBinding_Missing(t *testing.T) {
	s := &stubBinding{}
	b := bindings.CacheBinding{Delegate: s}

	if _, ok := b.GetAsBytes("test-unknown-key"); ok {
		t.Errorf("does not identify invalid key")
	}
	if _, ok := b.GetAsBytes("test-unknown-key"); ok {
		t.Errorf("does not identify invalid key")
	}
	if s.getAsBytesCount != 2 {
		t.Errorf("does not call delegate enough")
	}
}

func TestCacheBinding_Cached(t *testing.T) {
	s := &stubBinding{}
	b := bindings.CacheBinding{Delegate: s}

	if v, ok := b.GetAsBytes("test-key"); !ok || v == nil {
		t.Errorf("did not retrieve value as bytes")
	}
	if v, ok := b.GetAsBytes("test-key"); !ok || v == nil {
		t.Errorf("did not retrieve value as bytes")
	}
	if s.getAsBytesCount != 1 {
		t.Errorf("did not call delegate enough")
	}
}

func TestCacheBinding_GetName(t *testing.T) {
	s := &stubBinding{}
	b := bindings.CacheBinding{Delegate: s}

	if b.GetName() == "" {
		t.Errorf("did not retrieve name")
	}
	if b.GetName() == "" {
		t.Errorf("did not retrieve name")
	}
	if s.getNameCount != 2 {
		t.Errorf("did not call delegate enough")
	}
}

func TestConfigTreeBinding_GetAsBytes_MissingKey(t *testing.T) {
	b := bindings.ConfigTreeBinding{
		Root: filepath.Join("testdata", "test-k8s"),
	}

	if _, ok := b.GetAsBytes("test-missing-key"); ok {
		t.Errorf("does not identify missing key")
	}
}

func TestConfigTreeBinding_GetAsBytes_Directory(t *testing.T) {
	b := bindings.ConfigTreeBinding{
		Root: filepath.Join("testdata", "test-k8s"),
	}

	if _, ok := b.GetAsBytes(".hidden-data"); ok {
		t.Errorf("does not identify directory")
	}
}

func TestConfigTreeBinding_GetAsBytes_InvalidKey(t *testing.T) {
	b := bindings.ConfigTreeBinding{
		Root: filepath.Join("testdata", "test-k8s"),
	}

	if _, ok := b.GetAsBytes("test^secret^key"); ok {
		t.Errorf("does not identify invalid key")
	}
}

func TestConfigTreeBinding_GetAsBytes_ValidKey(t *testing.T) {
	b := bindings.ConfigTreeBinding{
		Root: filepath.Join("testdata", "test-k8s"),
	}

	if v, ok := b.GetAsBytes("test-secret-key"); !ok {
		t.Errorf("does not identify valid key")
	} else if !bytes.Equal([]byte("test-secret-value\n"), v) {
		t.Errorf("returned the wrong value")
	}
}

func TestConfigTreeBinding_GetName(t *testing.T) {
	b := bindings.ConfigTreeBinding{
		Root: filepath.Join("testdata", "test-k8s"),
	}

	if b.GetName() != "test-k8s" {
		t.Errorf("returned the wrong value")
	}
}

func TestMapBinding_GetAsBytes_MissingKey(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"test-secret-key": []byte("test-secret-value\n"),
		},
	}

	if _, ok := b.GetAsBytes("test-missing-key"); ok {
		t.Errorf("does not identify missing key")
	}
}

func TestMapBinding_GetAsBytes_InvalidKey(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"test-secret-key": []byte("test-secret-value\n"),
		},
	}

	if _, ok := b.GetAsBytes("test^secret^key"); ok {
		t.Errorf("does not identify invalid key")
	}
}

func TestMapBinding_GetAsBytes_ValidKey(t *testing.T) {
	b := bindings.MapBinding{
		Name: "test-name",
		Content: map[string][]byte{
			"test-secret-key": []byte("test-secret-value\n"),
		},
	}

	if v, ok := b.GetAsBytes("test-secret-key"); !ok {
		t.Errorf("does not identify valid key")
	} else if !bytes.Equal([]byte("test-secret-value\n"), v) {
		t.Errorf("returned the wrong value")
	}
}

func TestMapBinding_GetName(t *testing.T) {
	b := bindings.MapBinding{
		Name:    "test-name",
		Content: map[string][]byte{},
	}

	if b.GetName() != "test-name" {
		t.Errorf("returned the wrong value")
	}
}

type stubBinding struct {
	getAsBytesCount int
	getNameCount    int
}

func (s *stubBinding) GetAsBytes(key string) ([]byte, bool) {
	s.getAsBytesCount++

	if key == "test-key" {
		return []byte{}, true
	}

	return nil, false
}

func (s *stubBinding) GetName() string {
	s.getNameCount++
	return "test-name"
}
