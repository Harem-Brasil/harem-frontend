package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harem-brasil/backend/internal/domain"
)

// BindStrictJSON decodes JSON from the request body into dest, rejecting unknown fields.
func BindStrictJSON(c *gin.Context, dest any) error {
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}

// IsSyntaxOrUnknownField reports whether a JSON decode error is due to malformed
// syntax, an unknown field, an unexpected EOF, or a type mismatch.
func IsSyntaxOrUnknownField(err error) bool {
	if err == nil {
		return false
	}
	var synErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	return errors.As(err, &synErr) || errors.As(err, &typeErr) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "unknown field")
}

// RespondProblem envia application/problem+json (RFC 7807).
func RespondProblem(c *gin.Context, status int, title, detail string) {
	c.Header("Content-Type", "application/problem+json; charset=utf-8")
	c.JSON(status, domain.ProblemDetail{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

// RespondJSON define Content-Type JSON explícito (rotas API).
func RespondJSON(c *gin.Context, status int, data any) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(status, data)
}

// RespondValidation envia 422 com campos inválidos.
func RespondValidation(c *gin.Context, detail string, fields map[string]string) {
	c.Header("Content-Type", "application/problem+json; charset=utf-8")
	c.JSON(http.StatusUnprocessableEntity, domain.ProblemDetail{
		Type:       "validation-error",
		Title:      "Validation Error",
		Status:     http.StatusUnprocessableEntity,
		Detail:     detail,
		Extensions: map[string]any{"fields": fields},
	})
}

// HandleServiceError mapeia domain.AppError e erros genéricos para resposta HTTP.
func HandleServiceError(c *gin.Context, logger *slog.Logger, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		if len(appErr.FieldErrors) > 0 {
			RespondValidation(c, appErr.Detail, appErr.FieldErrors)
			return
		}
		title := appErr.Title
		if title == "" {
			title = http.StatusText(appErr.HTTPStatus)
		}
		RespondProblem(c, appErr.HTTPStatus, title, appErr.Detail)
		return
	}
	if logger != nil {
		logger.Error("internal server error", "path", c.FullPath(), "error", err)
	}
	RespondProblem(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError),
		"An unexpected error occurred")
}
