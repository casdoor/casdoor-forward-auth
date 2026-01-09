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

package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	config := Config{
		CasdoorEndpoint:     "http://test.example.com",
		CasdoorClientId:     "test-client-id",
		CasdoorClientSecret: "test-client-secret",
		CasdoorOrganization: "test-org",
		CasdoorApplication:  "test-app",
		PluginEndpoint:      "http://localhost:9999",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	tmpfile.Close()

	// Load the config
	LoadConfigFile(tmpfile.Name())

	// Verify the config was loaded correctly
	if CurrentConfig.CasdoorEndpoint != "http://test.example.com" {
		t.Errorf("Expected endpoint http://test.example.com, got %s", CurrentConfig.CasdoorEndpoint)
	}

	if CurrentConfig.CasdoorClientId != "test-client-id" {
		t.Errorf("Expected client ID test-client-id, got %s", CurrentConfig.CasdoorClientId)
	}

	if CurrentConfig.CasdoorClientSecret != "test-client-secret" {
		t.Errorf("Expected client secret test-client-secret, got %s", CurrentConfig.CasdoorClientSecret)
	}

	if CurrentConfig.CasdoorOrganization != "test-org" {
		t.Errorf("Expected organization test-org, got %s", CurrentConfig.CasdoorOrganization)
	}

	if CurrentConfig.CasdoorApplication != "test-app" {
		t.Errorf("Expected application test-app, got %s", CurrentConfig.CasdoorApplication)
	}

	if CurrentConfig.PluginEndpoint != "http://localhost:9999" {
		t.Errorf("Expected plugin endpoint http://localhost:9999, got %s", CurrentConfig.PluginEndpoint)
	}
}
