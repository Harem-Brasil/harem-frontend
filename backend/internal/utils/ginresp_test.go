package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBindStrictJSON_UnknownFieldsRejected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Whitelist struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "valid fields only",
			body:    `{"email":"a@b.com","password":"secret"}`,
			wantErr: false,
		},
		{
			name:    "unknown field rejected",
			body:    `{"email":"a@b.com","password":"secret","role":"admin"}`,
			wantErr: true,
		},
		{
			name:    "multiple unknown fields rejected",
			body:    `{"email":"a@b.com","password":"secret","role":"admin","extra":"x"}`,
			wantErr: true,
		},
		{
			name:    "nested unknown field rejected",
			body:    `{"email":"a@b.com","password":"secret","nested":{"foo":"bar"}}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON syntax",
			body:    `{"email":"a@b.com",`,
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(tt.body)))
			c.Request.Header.Set("Content-Type", "application/json")

			var dest Whitelist
			err := BindStrictJSON(c, &dest)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if dest.Email == "" {
					t.Fatalf("expected email to be parsed")
				}
			}
		})
	}
}
