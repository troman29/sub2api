//go:build unit

package service

import "testing"

func TestIsThinkingBlockSignatureErrorRecognizesInvalidEncryptedContent(t *testing.T) {
	service := &GatewayService{}
	tests := []struct {
		name string
		body string
	}{
		{
			name: "stable error code",
			body: `{"error":{"code":"invalid_encrypted_content","message":"opaque upstream rejection"}}`,
		},
		{
			name: "verification failed",
			body: `{"error":{"code":"invalid_encrypted_content","message":"The encrypted content could not be verified."}}`,
		},
		{
			name: "decryption or parsing failed",
			body: `{"error":{"code":"invalid_encrypted_content","message":"Encrypted content could not be decrypted or parsed."}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !service.isThinkingBlockSignatureError([]byte(test.body)) {
				t.Fatal("expected invalid encrypted content to trigger thinking-block rectification")
			}
		})
	}
}

func TestIsThinkingBlockSignatureErrorIgnoresUnrelatedEncryptedContent(t *testing.T) {
	service := &GatewayService{}
	body := []byte(`{"error":{"message":"Encrypted content is temporarily unavailable."}}`)

	if service.isThinkingBlockSignatureError(body) {
		t.Fatal("unrelated encrypted-content errors must not trigger rectification")
	}
}

func TestRectifierModelForEncryptedContentUsesOriginalClientModel(t *testing.T) {
	service := &GatewayService{}
	body := []byte(`{"error":{"message":"The encrypted content could not be verified."}}`)

	model := service.rectifierModelForError(body, "gpt-5.6-terra", "claude-opus-5")
	if model != "claude-opus-5" {
		t.Fatalf("expected original Claude model, got %q", model)
	}
}

func TestRectifierModelForSignatureErrorUsesMappedModel(t *testing.T) {
	service := &GatewayService{}
	body := []byte(`{"error":{"message":"Invalid signature in thinking block."}}`)

	model := service.rectifierModelForError(body, "claude-opus-5", "claude-sonnet-5")
	if model != "claude-opus-5" {
		t.Fatalf("expected mapped model, got %q", model)
	}
}

func TestRectifierModelForEncryptedContentPreservesPassbackModel(t *testing.T) {
	service := &GatewayService{}
	body := []byte(`{"error":{"code":"invalid_encrypted_content","message":"opaque rejection"}}`)

	model := service.rectifierModelForError(body, "deepseek-v4", "deepseek-v4")
	if model != "deepseek-v4" {
		t.Fatalf("expected passback model to remain unchanged, got %q", model)
	}
}
