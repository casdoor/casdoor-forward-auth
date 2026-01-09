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

func TestNewStateMemoryStorage(t *testing.T) {
	storage, err := NewStateMemoryStorage()
	if err != nil {
		t.Fatalf("Failed to create StateMemoryStorage: %v", err)
	}
	
	if storage == nil {
		t.Fatal("Expected non-nil storage")
	}
	
	if storage.content == nil {
		t.Fatal("Expected non-nil content map")
	}
}

func TestSetState(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	state := &State{
		Method: "GET",
		Header: http.Header{},
		Body:   []byte("test"),
	}
	
	nonce, err := storage.SetState(state)
	if err != nil {
		t.Fatalf("Failed to set state: %v", err)
	}
	
	if nonce == 0 {
		t.Error("Expected non-zero nonce")
	}
	
	// Verify state was stored
	retrievedState, err := storage.GetState(nonce)
	if err != nil {
		t.Fatalf("Failed to get state: %v", err)
	}
	
	if retrievedState.Method != "GET" {
		t.Errorf("Expected method GET, got %s", retrievedState.Method)
	}
}

func TestGetState(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	state := &State{
		Method: "POST",
		Header: http.Header{},
		Body:   []byte("test body"),
	}
	
	nonce, _ := storage.SetState(state)
	
	retrievedState, err := storage.GetState(nonce)
	if err != nil {
		t.Fatalf("Failed to get state: %v", err)
	}
	
	if retrievedState.Method != "POST" {
		t.Errorf("Expected method POST, got %s", retrievedState.Method)
	}
	
	if string(retrievedState.Body) != "test body" {
		t.Errorf("Expected body 'test body', got %s", string(retrievedState.Body))
	}
}

func TestGetStateNotFound(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	_, err := storage.GetState(12345)
	if err == nil {
		t.Error("Expected error when getting non-existent state")
	}
}

func TestPopState(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	state := &State{
		Method: "DELETE",
		Header: http.Header{},
		Body:   []byte("test"),
	}
	
	nonce, _ := storage.SetState(state)
	
	retrievedState, err := storage.PopState(nonce)
	if err != nil {
		t.Fatalf("Failed to pop state: %v", err)
	}
	
	if retrievedState.Method != "DELETE" {
		t.Errorf("Expected method DELETE, got %s", retrievedState.Method)
	}
	
	// Verify state was removed
	_, err = storage.GetState(nonce)
	if err == nil {
		t.Error("Expected error when getting popped state")
	}
}

func TestPopStateNotFound(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	_, err := storage.PopState(99999)
	if err == nil {
		t.Error("Expected error when popping non-existent state")
	}
}

func TestConcurrentAccess(t *testing.T) {
	storage, _ := NewStateMemoryStorage()
	
	done := make(chan bool)
	
	// Test concurrent SetState
	for i := 0; i < 10; i++ {
		go func(i int) {
			state := &State{
				Method: "GET",
				Header: http.Header{},
				Body:   []byte("test"),
			}
			_, err := storage.SetState(state)
			if err != nil {
				t.Errorf("Concurrent SetState failed: %v", err)
			}
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
}
