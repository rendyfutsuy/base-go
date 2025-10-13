package unittest

import (
	"testing"

	"github.com/rendyfutsuy/base-go/utils"
)

func TestIsDocumentNumberRisk(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		expected       bool
	}{
		{
			name:           "Valid risk document number",
			documentNumber: "FAC12345-R1",
			expected:       true,
		},
		{
			name:           "Valid risk document number with multiple digits",
			documentNumber: "FAC987654321-R123",
			expected:       true,
		},
		{
			name:           "Missing -R part",
			documentNumber: "FAC12345",
			expected:       false,
		},
		{
			name:           "Missing FAC prefix",
			documentNumber: "12345-R1",
			expected:       false,
		},
		{
			name:           "Extra characters at end",
			documentNumber: "FAC12345-R1X",
			expected:       false,
		},
		{
			name:           "Lowercase fac",
			documentNumber: "fac12345-R1",
			expected:       false,
		},
		{
			name:           "No digits after FAC",
			documentNumber: "FAC-R1",
			expected:       false,
		},
		{
			name:           "No digits after -R",
			documentNumber: "FAC12345-R",
			expected:       false,
		},
		{
			name:           "Empty string",
			documentNumber: "",
			expected:       false,
		},
		{
			name:           "add -RL",
			documentNumber: "FAC12345-RL1",
			expected:       false,
		},
		{
			name:           "add -R1-RL",
			documentNumber: "FAC12345-R1-RL1",
			expected:       false,
		},
		{
			name:           "add -R1-E1",
			documentNumber: "FAC12345-R1-E1",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsDocumentNumberRisk(tt.documentNumber)
			if result != tt.expected {
				t.Errorf("IsDocumentNumberRisk(%q) = %v; want %v", tt.documentNumber, result, tt.expected)
			}
		})
	}
}
