package dify

import (
	"testing"

	"github.com/openexw/dify-go"
	"github.com/stretchr/testify/assert"
)

// TestNewClient tests the creation of a new Dify client
func TestNewClient(t *testing.T) {
	// Create a new client
	client := dify.NewClient()

	// Verify the client is not nil
	assert.NotNil(t, client, "Client should not be nil")

	// Verify the Workflow method returns an instance
	appKey := "test-app-key"
	workflow := client.Workflow(appKey)
	assert.NotNil(t, workflow, "Workflow instance should not be nil")
}

// TestClient_Workflow tests the Workflow method with different appKeys
func TestClient_Workflow(t *testing.T) {
	// Create client
	client := dify.NewClient()

	// Test cases with different appKey scenarios
	testCases := []struct {
		name   string
		appKey string
	}{{
		name:   "Valid appKey",
		appKey: "valid-app-key",
	}, {
		name:   "Empty appKey",
		appKey: "",
	}, {
		name:   "Long appKey",
		appKey: "very-long-app-key-that-should-still-work",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			workflow := client.Workflow(tc.appKey)
			assert.NotNil(t, workflow, "Workflow instance should not be nil")
		})
	}
}
