package errorResponse

import (
	"encoding/json"
	"errors"
	"link_watcher/serviceErrors"
	"net/http"
)

type ErrorCode string

const (
	CodeNotFound   ErrorCode = "NOT_FOUND"
	CodeInternal   ErrorCode = "INTERNAL_ERROR"
	CodeConflict   ErrorCode = "CONFLICT"
	CodeBadRequest ErrorCode = "BAD_REQUEST"
)

type UserMessage struct {
	LangEn string `json:"langEn"`
	LangRu string `json:"langRu"`
}

type ErrorResponse struct {
	Code         ErrorCode   `json:"code"`
	InternalCode string      `json:"internalCode"`
	DevMessage   string      `json:"devMessage"`
	UserMessage  UserMessage `json:"userMessage"`
}

type errorResponse struct {
	StatusCode   int
	Code         ErrorCode
	InternalCode string
	DevMessage   string
	UserMessage  UserMessage
}

var errorResponseMessage = map[error]errorResponse{
	serviceErrors.ErrBadRequest: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_REQUEST_DATA",
		DevMessage:   "Invalid Request Data",
		UserMessage: UserMessage{
			LangEn: "Invalid Request Data",
			LangRu: "Некорректные данные запроса",
		},
	},
	serviceErrors.ErrInvalidId: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_ID",
		DevMessage:   "Invalid Id format",
		UserMessage: UserMessage{
			LangEn: "Invalid Id format",
			LangRu: "Некорректный формат Id",
		},
	},
	serviceErrors.ErrInvalidJSON: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_JSON",
		DevMessage:   "Invalid JSON body",
		UserMessage: UserMessage{
			LangEn: "Invalid JSON body",
			LangRu: "Некорректное тело запроса",
		},
	},
	serviceErrors.ErrEmptyURL: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "EMPTY_URL",
		DevMessage:   "URL is empty",
		UserMessage: UserMessage{
			LangEn: "URL is required",
			LangRu: "URL обязателен",
		},
	},
	serviceErrors.ErrInvalidURL: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_URL",
		DevMessage:   "Invalid URL",
		UserMessage: UserMessage{
			LangEn: "Invalid URL",
			LangRu: "Некорректный URL",
		},
	},
	serviceErrors.ErrURLTooLarge: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "URL_TOO_LARGE",
		DevMessage:   "URL is too large",
		UserMessage: UserMessage{
			LangEn: "URL is too large",
			LangRu: "URL слишком длинный",
		},
	},
	serviceErrors.ErrInvalidInterval: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_INTERVAL",
		DevMessage:   "Interval error",
		UserMessage: UserMessage{
			LangEn: "Interval must be between 1 and 60 seconds",
			LangRu: "Интервал должен быть от 1 до 60 секунд",
		},
	},
	serviceErrors.ErrNotFound: {
		StatusCode:   http.StatusNotFound,
		Code:         CodeNotFound,
		InternalCode: "TARGET_NOT_FOUND",
		DevMessage:   "Target not found",
		UserMessage: UserMessage{
			LangEn: "Target not found",
			LangRu: "Ссылка не найдена",
		},
	},
	serviceErrors.ErrConflict: {
		StatusCode:   http.StatusConflict,
		Code:         CodeConflict,
		InternalCode: "TARGET_CONFLICT",
		DevMessage:   "Target already exists",
		UserMessage: UserMessage{
			LangEn: "Target already exists",
			LangRu: "Ссылка уже отслеживается",
		},
	},
	serviceErrors.ErrInternal: {
		StatusCode:   http.StatusInternalServerError,
		Code:         CodeInternal,
		InternalCode: "INTERNAL_ERROR",
		DevMessage:   "Internal Error",
		UserMessage: UserMessage{
			LangEn: "Internal Error",
			LangRu: "Внутренняя ошибка сервера",
		},
	},
}

func ErrorResponseJSON(w http.ResponseWriter, errs error) {
	info, ok := errorResponseMessage[serviceErrors.ErrInternal]
	for target, mapped := range errorResponseMessage {
		if errors.Is(errs, target) {
			info = mapped
			ok = true
			break
		}
	}
	if !ok {
		info = errorResponseMessage[serviceErrors.ErrInternal]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(info.StatusCode)

	resp := ErrorResponse{
		Code:         info.Code,
		InternalCode: info.InternalCode,
		DevMessage:   info.DevMessage,
		UserMessage:  info.UserMessage,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(`{
		"code":"INTERNAL",
		"internalCode":"MARSHAL_ERROR",
		"devMessage":"failed to marshal error",
		"userMessage":{"langEn":"Internal server error","langRu":"Внутренняя ошибка сервера"}}`))
	}

	w.Write(body)

}
