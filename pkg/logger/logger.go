package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config define a configuração do logger.
type Config struct {
	// Level define o nível mínimo de log (debug, info, warn, error, fatal)
	Level string
	// Environment define o ambiente (dev, staging, production)
	Environment string
	// Encoding define o formato de saída (json, console)
	Encoding string
	// OutputPaths define os destinos de saída dos logs
	OutputPaths []string
	// ErrorOutputPaths define os destinos de saída dos logs de erro
	ErrorOutputPaths []string
	// EnableSampling ativa o sampling para alta performance
	EnableSampling bool
	// SamplingInitial define quantos logs iguais serão mantidos inicialmente
	SamplingInitial int
	// SamplingThereafter define quantos logs iguais serão mantidos depois
	SamplingThereafter int
	// GlobalFields define campos que serão adicionados a todos os logs
	GlobalFields map[string]interface{}
}

// DefaultConfig retorna uma configuração padrão baseada no ambiente.
func DefaultConfig(environment string) Config {
	env := strings.ToLower(strings.TrimSpace(environment))

	config := Config{
		Environment:        env,
		Encoding:           "json",
		OutputPaths:        []string{"stdout"},
		ErrorOutputPaths:   []string{"stderr"},
		EnableSampling:     true,
		SamplingInitial:    100,
		SamplingThereafter: 100,
		GlobalFields:       make(map[string]interface{}),
	}

	switch env {
	case "production", "prod":
		config.Level = "info"
		config.Environment = "production"
	case "staging", "stage":
		config.Level = "info"
		config.Environment = "staging"
	default:
		config.Level = "debug"
		config.Environment = "development"
		config.Encoding = "console"
	}

	// Adicionar hostname como campo global
	if hostname, err := os.Hostname(); err == nil {
		config.GlobalFields["hostname"] = hostname
	}

	return config
}

// New cria um novo logger com a configuração fornecida.
func New(config Config) (*zap.Logger, error) {
	// Converter string de nível para zapcore.Level
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		level = zap.InfoLevel
	}

	// Configuração do encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Para console em dev, adicionar cores
	if config.Encoding == "console" && config.Environment == "development" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Criar configuração base
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       config.Environment == "development",
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          config.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       config.OutputPaths,
		ErrorOutputPaths:  config.ErrorOutputPaths,
	}

	// Configurar sampling apenas em produção
	if config.EnableSampling && config.Environment == "production" {
		zapConfig.Sampling = &zap.SamplingConfig{
			Initial:    config.SamplingInitial,
			Thereafter: config.SamplingThereafter,
		}
	}

	// Construir logger
	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1), // Skip para mostrar o caller correto
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	// Adicionar campos globais
	if len(config.GlobalFields) > 0 {
		fields := make([]zap.Field, 0, len(config.GlobalFields)+1)

		// Sempre adicionar environment
		fields = append(fields, zap.String("environment", config.Environment))

		for key, value := range config.GlobalFields {
			fields = append(fields, zap.Any(key, value))
		}

		logger = logger.With(fields...)
	} else {
		logger = logger.With(zap.String("environment", config.Environment))
	}

	return logger, nil
}

// NewDevelopment cria um logger otimizado para desenvolvimento.
func NewDevelopment() *zap.Logger {
	config := DefaultConfig("development")

	// Sobrescrever algumas configurações para dev
	config.Level = "debug"
	config.Encoding = "console"
	config.EnableSampling = false

	logger, err := New(config)
	if err != nil {
		panic("failed to create development logger: " + err.Error())
	}

	return logger
}

// NewProduction cria um logger otimizado para produção.
func NewProduction() *zap.Logger {
	config := DefaultConfig("production")

	// Sobrescrever algumas configurações para prod
	config.Level = "info"
	config.Encoding = "json"
	config.EnableSampling = true

	logger, err := New(config)
	if err != nil {
		panic("failed to create production logger: " + err.Error())
	}

	return logger
}

// Setup inicializa o logger global do projeto (mantido para compatibilidade).
// Deprecated: Use New(), NewDevelopment() ou NewProduction() e propague via context.
func Setup(mode string) func() {
	var (
		cfg    zap.Config
		logger *zap.Logger
		err    error
	)

	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "release", "production":
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

	undo := zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return func() {
		_ = logger.Sync()
		undo()
	}
}

// WithAppInfo adiciona informações da aplicação ao logger.
func WithAppInfo(logger *zap.Logger, appName, version string) *zap.Logger {
	return logger.With(
		zap.String("app_name", appName),
		zap.String("version", version),
	)
}
