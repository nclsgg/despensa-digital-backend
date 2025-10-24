package logger

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

// responseWriter é um wrapper do http.ResponseWriter que captura o status code.
type responseWriter struct {
	http.ResponseWriter
	status       int
	written      int64
	wroteHeader  bool
	captureError bool
	errorMsg     string
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Hijack implementa http.Hijacker
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// Flush implementa http.Flusher
func (rw *responseWriter) Flush() {
	f, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// LoggingMiddleware é um middleware HTTP que loga todas as requisições.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Criar contexto com request ID e logger
			ctx := r.Context()
			ctx = WithRequestID(ctx)
			ctx = WithLogger(ctx, logger)

			// Obter request ID do contexto
			requestID := GetRequestID(ctx)

			// Criar logger com informações da requisição
			reqLogger := logger.With(
				zap.String(FieldRequestID, requestID),
				zap.String(FieldHTTPMethod, r.Method),
				zap.String(FieldHTTPPath, r.URL.Path),
				zap.String(FieldHTTPRemoteAddr, r.RemoteAddr),
			)

			// Atualizar contexto com logger enriquecido
			ctx = WithLogger(ctx, reqLogger)
			r = r.WithContext(ctx)

			// Wrapper do ResponseWriter
			rw := newResponseWriter(w)

			// Capturar tempo de início
			start := time.Now()

			// Log de início da requisição
			reqLogger.Info("request started",
				zap.String(FieldHTTPHost, r.Host),
				zap.String(FieldHTTPUserAgent, r.UserAgent()),
				zap.String(FieldHTTPScheme, getScheme(r)),
			)

			// Recuperar de panics
			defer func() {
				if err := recover(); err != nil {
					// Log do panic
					reqLogger.Error("panic recovered",
						zap.Any("panic", err),
						zap.String(FieldStackTrace, string(debug.Stack())),
						zap.Int(FieldHTTPStatus, http.StatusInternalServerError),
					)

					// Retornar erro 500
					if !rw.wroteHeader {
						rw.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()

			// Executar próximo handler
			next.ServeHTTP(rw, r)

			// Calcular duração
			duration := time.Since(start)

			// Determinar nível de log baseado no status
			logFunc := reqLogger.Info
			if rw.status >= 500 {
				logFunc = reqLogger.Error
			} else if rw.status >= 400 {
				logFunc = reqLogger.Warn
			}

			// Log de fim da requisição
			logFunc("request completed",
				zap.Int(FieldHTTPStatus, rw.status),
				zap.Int64(FieldDuration, duration.Milliseconds()),
				zap.Int64(FieldHTTPResponseSize, rw.written),
			)
		})
	}
}

// GinLoggingMiddleware é um middleware compatível com Gin que loga todas as requisições.
func GinLoggingMiddleware(logger *zap.Logger) func(c interface{}) {
	// Usamos interface{} para evitar dependência circular do Gin
	// Na prática, c será *gin.Context
	return func(c interface{}) {
		// Type assertion para gin.Context seria feita aqui
		// mas mantemos genérico para este exemplo

		// Este é um placeholder - a implementação real precisaria do tipo gin.Context
		// Exemplo de uso seria:
		// func GinLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
		//     return func(c *gin.Context) {
		//         start := time.Now()
		//         ... implementação específica do Gin
		//     }
		// }
	}
}

// RecoveryMiddleware é um middleware que recupera de panics e loga o erro.
func RecoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Obter logger do contexto ou usar o fornecido
					log := FromContext(r.Context())
					if log == nil {
						log = logger
					}

					// Log do panic
					log.Error("panic recovered",
						zap.Any("panic", err),
						zap.String(FieldStackTrace, string(debug.Stack())),
						zap.String(FieldHTTPMethod, r.Method),
						zap.String(FieldHTTPPath, r.URL.Path),
					)

					// Retornar erro 500
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adiciona um request ID único a cada requisição.
// Use este middleware antes do LoggingMiddleware se quiser controle separado.
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithRequestID(r.Context())
			requestID := GetRequestID(ctx)

			// Adicionar request ID no header da resposta
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getScheme determina o esquema HTTP (http ou https) da requisição.
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}

	// Verificar headers de proxy
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}

	return "http"
}

// SkipPathsMiddleware permite pular o logging de certos paths (útil para health checks).
func SkipPathsMiddleware(logger *zap.Logger, skipPaths []string) func(http.Handler) http.Handler {
	skipMap := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skipMap[path] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Se o path deve ser pulado, executar sem logging
			if skipMap[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			// Caso contrário, usar o middleware de logging
			LoggingMiddleware(logger)(next).ServeHTTP(w, r)
		})
	}
}

// LogOnlyErrorsMiddleware loga apenas requisições com erro (status >= 400).
func LogOnlyErrorsMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Criar contexto com request ID e logger
			ctx := WithRequestID(r.Context())
			ctx = WithLogger(ctx, logger)
			r = r.WithContext(ctx)

			// Wrapper do ResponseWriter
			rw := newResponseWriter(w)

			// Capturar tempo de início
			start := time.Now()

			// Executar próximo handler
			next.ServeHTTP(rw, r)

			// Log apenas se houver erro
			if rw.status >= 400 {
				duration := time.Since(start)

				logFunc := logger.Warn
				if rw.status >= 500 {
					logFunc = logger.Error
				}

				logFunc("request failed",
					zap.String(FieldRequestID, GetRequestID(ctx)),
					zap.String(FieldHTTPMethod, r.Method),
					zap.String(FieldHTTPPath, r.URL.Path),
					zap.Int(FieldHTTPStatus, rw.status),
					zap.Int64(FieldDuration, duration.Milliseconds()),
				)
			}
		})
	}
}
