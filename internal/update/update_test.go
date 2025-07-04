package update_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gypsydave5/pairstair/internal/update"
)

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		mockResponse   string
		mockStatusCode int
		expectMessage  bool
		expectedMsg    string
	}{
		{
			name:           "newer version available",
			currentVersion: "v0.5.0",
			mockResponse: `[
				{"tag_name": "v0.6.0", "draft": false},
				{"tag_name": "v0.5.0", "draft": false}
			]`,
			mockStatusCode: 200,
			expectMessage:  true,
			expectedMsg:    "A newer version of pairstair is available: v0.6.0 (you have v0.5.0)",
		},
		{
			name:           "current version is latest",
			currentVersion: "v0.6.0",
			mockResponse: `[
				{"tag_name": "v0.6.0", "draft": false},
				{"tag_name": "v0.5.0", "draft": false}
			]`,
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "development version newer than latest",
			currentVersion: "0.7.0-dev+abc1234",
			mockResponse: `[
				{"tag_name": "v0.6.0", "draft": false}
			]`,
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "API failure returns empty string",
			currentVersion: "v0.5.0",
			mockResponse:   "",
			mockStatusCode: 500,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "empty response returns empty string",
			currentVersion: "v0.5.0",
			mockResponse:   "[]",
			mockStatusCode: 200,
			expectMessage:  false,
			expectedMsg:    "",
		},
		{
			name:           "draft releases are ignored",
			currentVersion: "v0.5.0",
			mockResponse: `[
				{"tag_name": "v0.7.0", "draft": true},
				{"tag_name": "v0.6.0", "draft": false}
			]`,
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
				if tt.mockResponse != "" {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer server.Close()

			// Test the update check function using the public API
			message := update.CheckForUpdateWithURL(tt.currentVersion, server.URL)

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

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name        string
		current     string
		latest      string
		expectNewer bool
	}{
		{
			name:        "newer version available",
			current:     "v0.5.0",
			latest:      "v0.6.0",
			expectNewer: true,
		},
		{
			name:        "same version",
			current:     "v0.5.0",
			latest:      "v0.5.0",
			expectNewer: false,
		},
		{
			name:        "current is newer",
			current:     "v0.6.0",
			latest:      "v0.5.0",
			expectNewer: false,
		},
		{
			name:        "development version vs release",
			current:     "0.5.0-dev+abc1234",
			latest:      "v0.6.0",
			expectNewer: true,
		},
		{
			name:        "development version newer than release",
			current:     "0.7.0-dev+abc1234",
			latest:      "v0.6.0",
			expectNewer: false,
		},
		{
			name:        "handle missing v prefix",
			current:     "0.5.0",
			latest:      "v0.6.0",
			expectNewer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := update.IsNewerVersion(tt.current, tt.latest)
			if result != tt.expectNewer {
				t.Errorf("IsNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, result, tt.expectNewer)
			}
		})
	}
}
