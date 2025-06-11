package docs

import (
	"encoding/json"
	"testing"

	"github.com/swaggo/swag"
)

// TestSwaggerSpecProcessing verifies that the Swagger template has been correctly
// processed and that the final JSON document contains valid data.
func TestSwaggerSpecProcessing(t *testing.T) {
	// Step 1: Get the processed document registered by the init() function
	specString, err := swag.ReadDoc("swagger")
	if err != nil {
		t.Fatalf("Cannot read the registered swagger specification: %v", err)
	}

	// Step 2: Parse the JSON document into a map to check its content
	var parsedSpec map[string]interface{}
	if err := json.Unmarshal([]byte(specString), &parsedSpec); err != nil {
		t.Fatalf("Failed to parse swagger specification as JSON: %v", err)
	}

	// Step 3: Check if the values in JSON match the source values in SwaggerInfo.
	// This confirms that the template {{.Title}} was correctly replaced.
	infoBlock, ok := parsedSpec["info"].(map[string]interface{})
	if !ok {
		t.Fatal("The 'info' block was not found in the swagger specification")
	}

	title, ok := infoBlock["title"].(string)
	if !ok {
		t.Fatal("The 'title' field was not found in the 'info' block")
	}

	if title != SwaggerInfo.Title {
		t.Errorf("expected title '%s', got '%s'", SwaggerInfo.Title, title)
	}

	// Step 4: Check another value for certainty, e.g., Host
	host, ok := parsedSpec["host"].(string)
	if !ok {
		t.Fatal("The 'host' field was not found in the swagger specification")
	}

	if host != SwaggerInfo.Host {
		t.Errorf("expected host '%s', got '%s'", SwaggerInfo.Host, host)
	}
}
