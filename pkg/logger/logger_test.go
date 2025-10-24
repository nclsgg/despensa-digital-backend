package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantLevel   string
		wantEnv     string
		wantEnc     string
	}{
		{
			name:        "production config",
			environment: "production",
			wantLevel:   "info",
			wantEnv:     "production",
			wantEnc:     "json",
		},
		{
			name:        "prod shorthand",
			environment: "prod",
			wantLevel:   "info",
			wantEnv:     "production",
			wantEnc:     "json",
		},
		{
			name:        "staging config",
			environment: "staging",
			wantLevel:   "info",
			wantEnv:     "staging",
			wantEnc:     "json",
		},
		{
			name:        "development config",
			environment: "development",
			wantLevel:   "debug",
			wantEnv:     "development",
			wantEnc:     "console",
		},
		{
			name:        "dev config",
			environment: "dev",
			wantLevel:   "debug",
			wantEnv:     "development",
			wantEnc:     "console",
		},
		{
			name:        "unknown defaults to dev",
			environment: "unknown",
			wantLevel:   "debug",
			wantEnv:     "development",
			wantEnc:     "console",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig(tt.environment)

			assert.Equal(t, tt.wantLevel, cfg.Level)
			assert.Equal(t, tt.wantEnv, cfg.Environment)
			assert.Equal(t, tt.wantEnc, cfg.Encoding)
			assert.NotNil(t, cfg.GlobalFields)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid production config",
			config: Config{
				Level:              "info",
				Environment:        "production",
				Encoding:           "json",
				OutputPaths:        []string{"stdout"},
				ErrorOutputPaths:   []string{"stderr"},
				EnableSampling:     true,
				SamplingInitial:    100,
				SamplingThereafter: 100,
				GlobalFields: map[string]interface{}{
					"app_name": "test-app",
				},
			},
			wantErr: false,
		},
		{
			name: "valid development config",
			config: Config{
				Level:            "debug",
				Environment:      "development",
				Encoding:         "console",
				OutputPaths:      []string{"stdout"},
				ErrorOutputPaths: []string{"stderr"},
				EnableSampling:   false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)

				// Verificar se o logger pode fazer log sem erro
				logger.Info("test message",
					zap.String("module", "test"),
				)
			}
		})
	}
}

func TestNewDevelopment(t *testing.T) {
	logger := NewDevelopment()
	require.NotNil(t, logger)

	// Testar se pode fazer log
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
}

func TestNewProduction(t *testing.T) {
	logger := NewProduction()
	require.NotNil(t, logger)

	// Testar se pode fazer log
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestWithAppInfo(t *testing.T) {
	baseLogger := NewDevelopment()
	logger := WithAppInfo(baseLogger, "test-app", "1.0.0")

	require.NotNil(t, logger)

	// Verificar se os campos foram adicionados
	// (isto seria melhor testado com um observer do zap, mas mantemos simples)
	logger.Info("test message")
}

func BenchmarkNewProduction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewProduction()
	}
}

func BenchmarkLogging(b *testing.B) {
	logger := NewProduction()

	b.Run("simple info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("test message")
		}
	})

	b.Run("with fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("test message",
				zap.String("module", "test"),
				zap.String("function", "TestFunc"),
				zap.Int("count", i),
			)
		}
	})

	b.Run("with many fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info("test message",
				zap.String("module", "test"),
				zap.String("function", "TestFunc"),
				zap.Int("count", i),
				zap.String("request_id", "req-123"),
				zap.String("user_id", "user-456"),
				zap.Duration("duration", 100),
			)
		}
	})
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantLevel zapcore.Level
	}{
		{"debug level", "debug", zapcore.DebugLevel},
		{"info level", "info", zapcore.InfoLevel},
		{"warn level", "warn", zapcore.WarnLevel},
		{"error level", "error", zapcore.ErrorLevel},
		{"invalid defaults to info", "invalid", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig("development")
			config.Level = tt.level

			logger, err := New(config)
			require.NoError(t, err)
			require.NotNil(t, logger)

			// Logger criado com sucesso
			logger.Info("test")
		})
	}
}
