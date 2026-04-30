package controllers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIsSyntaxOrUnknownField(t *testing.T) {
	tests := []struct {
		name  string
		input string
		into  any
		want  bool
	}{
		{
			name:  "unknown field",
			input: `{"email":"a@b.com","role":"admin"}`,
			into: &struct {
				Email string `json:"email"`
			}{},
			want: true,
		},
		{
			name:  "valid fields only",
			input: `{"email":"a@b.com"}`,
			into: &struct {
				Email string `json:"email"`
			}{},
			want: false,
		},
		{
			name:  "invalid JSON syntax",
			input: `{"email":`,
			into: &struct {
				Email string `json:"email"`
			}{},
			want: true,
		},
		{
			name:  "wrong type",
			input: `{"count":"not-a-number"}`,
			into: &struct {
				Count int `json:"count"`
			}{},
			want: true,
		},
		{
			name:  "empty string",
			input: ``,
			into: &struct {
				Email string `json:"email"`
			}{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := json.NewDecoder(strings.NewReader(tt.input))
			dec.DisallowUnknownFields()
			err := dec.Decode(tt.into)
			if err != nil {
				got := isSyntaxOrUnknownField(err)
				if got != tt.want {
					t.Errorf("isSyntaxOrUnknownField() = %v, want %v for error: %v", got, tt.want, err)
				}
			} else if tt.want {
				t.Fatalf("expected decode error")
			}
		})
	}
}
