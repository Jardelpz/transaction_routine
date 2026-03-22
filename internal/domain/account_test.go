package domain

import (
	"errors"
	"testing"
)

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name     string
		document string
		wantErr  error
	}{
		{
			name:     "valid document",
			document: "12345678901",
			wantErr:  nil,
		},
		{
			name:     "document too short",
			document: "1234567890",
			wantErr:  ErrDocumentInvalid,
		},
		{
			name:     "document too long",
			document: "123456789012",
			wantErr:  ErrDocumentInvalid,
		},
		{
			name:     "document with letters",
			document: "1234567890a",
			wantErr:  ErrDocumentNonNumber,
		},
		{
			name:     "document with special characters",
			document: "1234567890-",
			wantErr:  ErrDocumentNonNumber,
		},
		{
			name:     "empty document",
			document: "",
			wantErr:  ErrDocumentInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDocument(tt.document)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
