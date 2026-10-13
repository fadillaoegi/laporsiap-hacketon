package main

import "testing"

func TestDiagnosticExcludesRawProviderContent(t *testing.T) {
	data := []byte(`{"error":"rejected","provider_diagnostic":{"requested_model":"swiss-ai/apertus-v1.5-70b","response_model":"other/model","validation_rule":"PRIVATE-SOURCE","raw_body":"PRIVATE-SOURCE"}}`)
	d := safeDiagnostic(data)
	if d == nil || d.ResponseModel != "other/model" || d.Rule != "contract_validation_failed" {
		t.Fatal("safe diagnostic missing")
	}
	if safeDiagnostic([]byte(`{"provider_diagnostic":{"response_model":"private text with whitespace"}}`)) != nil {
		t.Fatal("unsafe model identifier retained")
	}
}
