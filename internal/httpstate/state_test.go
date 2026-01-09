// Copyright 2026 The Casdoor Authors. All Rights Reserved.
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

package httpstate

import (
	"net/http"
	"testing"
)

func TestNewState(t *testing.T) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("X-Custom-Header", "test-value")
	
	body := []byte("test body content")
	method := "POST"
	
	state := NewState(method, header, body)
	
	if state.Method != method {
		t.Errorf("Expected method %s, got %s", method, state.Method)
	}
	
	if state.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be application/json, got %s", state.Header.Get("Content-Type"))
	}
	
	if state.Header.Get("X-Custom-Header") != "test-value" {
		t.Errorf("Expected X-Custom-Header to be test-value, got %s", state.Header.Get("X-Custom-Header"))
	}
	
	if string(state.Body) != string(body) {
		t.Errorf("Expected body %s, got %s", string(body), string(state.Body))
	}
	
	// Verify deep copy of header
	header.Set("Content-Type", "text/html")
	if state.Header.Get("Content-Type") == "text/html" {
		t.Error("Header was not deep copied")
	}
	
	// Verify deep copy of body
	body[0] = 'X'
	if state.Body[0] == 'X' {
		t.Error("Body was not deep copied")
	}
}
