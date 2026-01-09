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

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Any("/test", TestHandler)
	
	body := []byte("test body")
	req, _ := http.NewRequest("POST", "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Header", "test-value")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	// The response should contain the modified body and headers
	if w.Body.String() == "" {
		t.Error("Expected non-empty response body")
	}
}

func TestReplacementStruct(t *testing.T) {
	replacement := Replacement{
		ShouldReplaceBody:   true,
		ShouldReplaceHeader: true,
		Body:                "test",
		Header:              make(map[string][]string),
	}
	
	if !replacement.ShouldReplaceBody {
		t.Error("Expected ShouldReplaceBody to be true")
	}
	
	if !replacement.ShouldReplaceHeader {
		t.Error("Expected ShouldReplaceHeader to be true")
	}
	
	if replacement.Body != "test" {
		t.Errorf("Expected body 'test', got %s", replacement.Body)
	}
}
