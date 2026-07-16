package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		expectedStatus int
		expectLog      bool
		logLevel       slog.Level
	}{
		{
			name: "successful request",
			handler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "server error",
			handler: func(c *gin.Context) {
				c.Status(http.StatusInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			logLevel:       slog.LevelError,
		},
		{
			name: "slow request",
			handler: func(c *gin.Context) {
				time.Sleep(1100 * time.Millisecond)
				c.Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectLog:      true,
			logLevel:       slog.LevelWarn,
		},
		{
			name: "client error",
			handler: func(c *gin.Context) {
				c.Status(http.StatusBadRequest)
			},
			expectedStatus: http.StatusBadRequest,
			expectLog:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			r := gin.New()
			r.Use(Logger(logger))
			r.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectLog {
				assert.NotEmpty(t, logBuffer.String())
				logOutput := logBuffer.String()
				switch tt.logLevel {
				case slog.LevelError:
					assert.Contains(t, logOutput, "request failed")
				case slog.LevelWarn:
					assert.Contains(t, logOutput, "slow request")
				}
			} else {
				assert.Empty(t, logBuffer.String())
			}
		})
	}
}

func TestLoggerWithClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	r := gin.New()
	r.Use(Logger(logger))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(w, req)

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "request failed")
	assert.Contains(t, logOutput, "192.168.1.1")
}

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		expectedStatus int
		expectLog      bool
		expectedLogMsg string
	}{
		{
			name: "no panic",
			handler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectLog:      false,
		},
		{
			name: "panic recovered",
			handler: func(c *gin.Context) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			expectedLogMsg: "panic recovered",
		},
		{
			name: "panic with string",
			handler: func(c *gin.Context) {
				panic("something went wrong")
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			expectedLogMsg: "panic recovered",
		},
		{
			name: "panic with int",
			handler: func(c *gin.Context) {
				panic(42)
			},
			expectedStatus: http.StatusInternalServerError,
			expectLog:      true,
			expectedLogMsg: "panic recovered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			r := gin.New()
			r.Use(Recovery(logger))
			r.GET("/test", tt.handler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectLog {
				assert.NotEmpty(t, logBuffer.String())
				logOutput := logBuffer.String()
				assert.Contains(t, logOutput, tt.expectedLogMsg)
				assert.Contains(t, logOutput, "method=GET")
				assert.Contains(t, logOutput, "path=/test")
			} else {
				assert.Empty(t, logBuffer.String())
			}
		})
	}
}

func TestRecoveryWithErrAbortHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("panic with http.ErrAbortHandler", func(t *testing.T) {
		var logBuffer bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		r := gin.New()
		r.Use(Recovery(logger))
		r.GET("/test", func(c *gin.Context) {
			panic(http.ErrAbortHandler)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		logOutput := logBuffer.String()
		assert.NotContains(t, logOutput, "panic recovered")
	})
}

func TestRecoveryResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	r := gin.New()
	r.Use(Recovery(logger))
	r.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])
	assert.Contains(t, logBuffer.String(), "panic recovered")
}

func TestLoggerAndRecoveryCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	r := gin.New()
	r.Use(Logger(logger))
	r.Use(Recovery(logger))
	r.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "request failed")
	assert.Contains(t, logOutput, "panic recovered")
}

func TestRecoveryWithDifferentPanicTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	panicTypes := []struct {
		name  string
		panic interface{}
	}{
		{"string panic", "test string"},
		{"error panic", http.ErrBodyNotAllowed},
		{"int panic", 123},
		{"bool panic", true},
		{"struct panic", struct{}{}},
	}

	for _, pt := range panicTypes {
		t.Run(pt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			r := gin.New()
			r.Use(Recovery(logger))
			r.GET("/test", func(c *gin.Context) {
				panic(pt.panic)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Contains(t, logBuffer.String(), "panic recovered")
		})
	}
}

func TestLoggerWithDifferentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))

			r := gin.New()
			r.Use(Logger(logger))
			r.Any("/test", func(c *gin.Context) {
				c.Status(http.StatusInternalServerError)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, "request failed")
			assert.Contains(t, logOutput, method)
		})
	}
}
