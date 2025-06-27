package main

import (
	"runtime/debug"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name        string
		buildInfo   *debug.BuildInfo
		expected    string
		description string
	}{
		{
			name: "fallback to constant when no build info",
			buildInfo: nil,
			expected: "0.6.0-dev",
			description: "Should return the version constant when debug.ReadBuildInfo() fails",
		},
		{
			name: "clean git tag",
			buildInfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.tag", Value: "v0.5.0"},
					{Key: "vcs.revision", Value: "abc123def456"},
					{Key: "vcs.modified", Value: "false"},
				},
			},
			expected: "v0.5.0",
			description: "Should return the git tag when available and not modified",
		},
		{
			name: "dirty git tag",
			buildInfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.tag", Value: "v0.5.0"},
					{Key: "vcs.revision", Value: "abc123def456"},
					{Key: "vcs.modified", Value: "true"},
				},
			},
			expected: "v0.5.0-dirty",
			description: "Should return tag with -dirty suffix when modified",
		},
		{
			name: "commit hash only, clean",
			buildInfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abc123def456789"},
					{Key: "vcs.modified", Value: "false"},
				},
			},
			expected: "0.6.0-dev+abc123de",
			description: "Should return version constant + short hash when no tag available",
		},
		{
			name: "commit hash only, dirty",
			buildInfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abc123def456789"},
					{Key: "vcs.modified", Value: "true"},
				},
			},
			expected: "0.6.0-dev+abc123de-dirty",
			description: "Should return version constant + short hash + dirty when no tag and modified",
		},
		{
			name: "short commit hash",
			buildInfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "abc123"},
					{Key: "vcs.modified", Value: "false"},
				},
			},
			expected: "0.6.0-dev+abc123",
			description: "Should handle short commit hashes without truncation",
		},
		{
			name: "module version",
			buildInfo: &debug.BuildInfo{
				Main: debug.Module{
					Version: "v0.4.1",
				},
				Settings: []debug.BuildSetting{},
			},
			expected: "v0.4.1",
			description: "Should return module version when available and no VCS info",
		},
		{
			name: "module version with devel",
			buildInfo: &debug.BuildInfo{
				Main: debug.Module{
					Version: "(devel)",
				},
				Settings: []debug.BuildSetting{},
			},
			expected: "0.6.0-dev",
			description: "Should fallback to constant when module version is (devel)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the actual production function with controlled input
			result := getVersionFromBuildInfo(tt.buildInfo, tt.buildInfo != nil)
			if result != tt.expected {
				t.Errorf("getVersionFromBuildInfo() = %q, want %q\nDescription: %s", 
					result, tt.expected, tt.description)
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	// Test that the version flag is properly defined
	config := parseFlags()
	
	// Check that Version field exists and defaults to false
	if config.Version != false {
		t.Errorf("Version flag should default to false, got %v", config.Version)
	}
}

func TestGetVersionWithNoBuildInfo(t *testing.T) {
	// Test the case where debug.ReadBuildInfo() fails (hasInfo = false)
	result := getVersionFromBuildInfo(nil, false)
	expected := "0.6.0-dev"
	
	if result != expected {
		t.Errorf("getVersionFromBuildInfo(nil, false) = %q, want %q", result, expected)
	}
}
