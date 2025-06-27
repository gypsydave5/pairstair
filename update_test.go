package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		mockResponse   []GitHubRelease
		mockStatusCode int
		expectMessage  bool
		expectedMsg    string
	}{
		{
			name:           "newer version available",
			currentVersion: "v0.5.0",
			mockResponse: []GitHubRelease{
				{TagName: "v0.6.0", Draft: false},
				{TagName: "v0.5.0", Draft: false},
			},
			mockStatusCode: 200,
			expectMessage:  true,
			expectedMsg:    "A newer version of pairstair is available: v0.6.0 (you have v0.5.0)",
		},
		{
			name:           "current version is latest",
			currentVersion: "v0.6.0",
			mockResponse: []GitHubRelease{
				{TagName: "v0.6.0", Draft: false},
				{TagName: "v0.5.0", Draft: false},
			},
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "development version newer than latest",
			currentVersion: "0.7.0-dev+abc1234",
			mockResponse: []GitHubRelease{
				{TagName: "v0.6.0", Draft: false},
			},
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "API failure returns empty string",
			currentVersion: "v0.5.0",
			mockResponse:   nil,
			mockStatusCode: 500,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "empty response returns empty string",
			currentVersion: "v0.5.0",
			mockResponse:   []GitHubRelease{},
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "draft releases are ignored",
			currentVersion: "v0.5.0",
			mockResponse: []GitHubRelease{
				{TagName: "v0.7.0", Draft: true},  // This should be ignored
				{TagName: "v0.6.0", Draft: false}, // This is the latest non-draft
			},
			mockStatusCode: 200,
			expectMessage:  true,
			expectedMsg:    "A newer version of pairstair is available: v0.6.0 (you have v0.5.0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			// Test the update check function
			message := checkForUpdateWithURL(tt.currentVersion, server.URL)

			if tt.expectMessage {
				if message == "" {
					t.Errorf("Expected update message, got empty string")
				}
				if message != tt.expectedMsg {
					t.Errorf("Expected message %q, got %q", tt.expectedMsg, message)
				}
			} else {
				if message != "" {
					t.Errorf("Expected no message, got %q", message)
				}
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name       string
		current    string
		latest     string
		expectNewer bool
	}{
		{
			name:       "newer version available",
			current:    "v0.5.0",
			latest:     "v0.6.0",
			expectNewer: true,
		},
		{
			name:       "same version",
			current:    "v0.5.0",
			latest:     "v0.5.0",
			expectNewer: false,
		},
		{
			name:       "current is newer",
			current:    "v0.6.0",
			latest:     "v0.5.0",
			expectNewer: false,
		},
		{
			name:       "development version vs release",
			current:    "0.5.0-dev+abc1234",
			latest:     "v0.6.0",
			expectNewer: true,
		},
		{
			name:       "development version newer than release",
			current:    "0.7.0-dev+abc1234",
			latest:     "v0.6.0",
			expectNewer: false,
		},
		{
			name:       "handle missing v prefix",
			current:    "0.5.0",
			latest:     "v0.6.0",
			expectNewer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNewerVersion(tt.current, tt.latest)
			if result != tt.expectNewer {
				t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, result, tt.expectNewer)
			}
		})
	}
}
