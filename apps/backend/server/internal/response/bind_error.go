package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"qingmo-server/internal/errs"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func parseBindError(err error) []FieldError {

	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		out := make([]FieldError, 0, len(ve))

		for _, fe := range ve {
			out = append(out, FieldError{
				Field:   jsonFieldName(fe),
				Message: fieldErrorMessage(fe),
			})
		}
		return out
	}
	// JSON 格式错误
	if syn, ok := errors.AsType[*json.SyntaxError](err); ok {
		return []FieldError{
			{
				Field:   "",
				Message: fmt.Sprintf("invalid JSON at offset %d", syn.Offset),
			},
		}
	}
	// JSON 解析类型错误
	if typ, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
		return []FieldError{
			{
				Field: typ.Field,
				Message: fmt.Sprintf(
					"expected %s, got %s",
					typ.Type,
					typ.Value,
				),
			},
		}
	}
	// 请求体为空
	if errors.Is(err, io.EOF) {
		return []FieldError{{Field: "", Message: "request body is empty"}}
	}

	return nil
}

func jsonFieldName(fe validator.FieldError) string {
	return toSnakeCase(fe.Field())
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "url":
		return "must be a valid URL"
	case "numeric":
		return "must be numeric"
	default:
		return fmt.Sprintf("failed on '%s'", fe.Tag())
	}
}

func BindError(c *gin.Context, err error) bool {
	fields := parseBindError(err)
	if fields != nil {
		ErrorWithData(
			c,
			http.StatusBadRequest,
			CodeInvalidParam,
			"invalid request body",
			fields,
		)
		return true
	}
	WriteError(c, errs.Wrap(errs.KindInvalid, "invalid request body", err))
	return true
}
