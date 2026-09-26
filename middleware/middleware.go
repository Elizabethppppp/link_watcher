package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ResponseLog struct {
	http.ResponseWriter
	size    int
	status  int
	headers http.Header
	body    []byte
}

func NewResponseLog(w http.ResponseWriter) *ResponseLog {
	return &ResponseLog{
		ResponseWriter: w,
		status:         http.StatusOK,
		headers:        make(http.Header),
		body:           make([]byte, 0),
	}
}

func (w *ResponseLog) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseLog) Write(p []byte) (int, error) {
	w.size += len(p)
	w.body = append(w.body, p...)
	return w.ResponseWriter.Write(p)
}

func (w *ResponseLog) Header() http.Header {
	return w.ResponseWriter.Header()
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var bodyBytes []byte
		if r.Body != nil {
			var err error
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				log.Printf("read request body: %v", err)
			}
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		requestFields := map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		}
		if len(r.Header) > 0 {
			requestFields["headers"] = r.Header
		}
		if len(bodyBytes) > 0 {
			requestFields["body"] = string(bodyBytes)
		}
		writeRequest("REQUEST", requestFields)

		rw := NewResponseLog(w)

		// Логируем ответ даже при панике в хендлере.
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic in handler: %v", rec)
				rw.status = http.StatusInternalServerError
			}

			responseFields := map[string]any{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rw.status,
				"size":        rw.size,
				"duration_ms": time.Since(start).Milliseconds(),
			}
			if len(rw.Header()) > 0 {
				responseFields["headers"] = rw.Header()
			}
			if body := string(rw.body); body != "" {
				if json.Valid(rw.body) {
					responseFields["body"] = json.RawMessage(rw.body)
				} else {
					responseFields["body"] = body
				}
			}
			writeResponse("RESPONSE", responseFields)
		}()

		next.ServeHTTP(rw, r)
	})
}

func writeRequest(label string, fields map[string]any) {
	fields["time"] = time.Now().Format(time.RFC3339)
	writeResponse(label, fields)
}

func writeResponse(label string, logMap map[string]interface{}) {
	jsonData, err := json.MarshalIndent(logMap, "", "  ")
	if err != nil {
		return
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return
	}

	file, err := os.OpenFile("logs/requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v", err)
		}
	}()
	separator := "--" + label + "--\n"
	_, _ = file.WriteString(separator)
	_, _ = file.Write(append(jsonData, '\n'))
}
