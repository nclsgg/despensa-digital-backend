package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nclsgg/despensa-digital/backend/pkg/logger"
	"go.uber.org/zap"
)

// Simula um modelo de negócio
type Order struct {
	ID       string
	UserID   string
	Status   string
	Items    int
	Total    float64
	CreateAt time.Time
}

// OrderService demonstra uso de logging em camada de serviço
type OrderService struct {
	// Em um app real, teria dependências como repository, etc
}

// CreateOrder processa um pedido com logging apropriado
func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
	// Obter logger do contexto e adicionar campos específicos
	log := logger.FromContext(ctx).With(
		zap.String(logger.FieldModule, "orders"),
		zap.String(logger.FieldFunction, "CreateOrder"),
		zap.String(logger.FieldEntityType, "order"),
		zap.String(logger.FieldEntityID, order.ID),
	)

	log.Info("starting order creation",
		zap.String("user_id", order.UserID),
		zap.Int("items_count", order.Items),
		zap.Float64("total", order.Total),
	)

	// Simular validação
	if order.Items == 0 {
		log.Warn("order validation failed",
			zap.String(logger.FieldReason, "no_items"),
		)
		return errors.New("order must have at least one item")
	}

	// Simular processamento com medição de tempo
	start := time.Now()

	// Simular algum processamento...
	time.Sleep(50 * time.Millisecond)

	duration := time.Since(start)

	// Simular possível erro
	if order.Total < 0 {
		log.Error("order processing failed",
			logger.WithError(logger.ErrorFields{
				Error:     errors.New("invalid total amount"),
				ErrorType: "ValidationError",
				ErrorCode: "ORD001",
				Operation: "validate_total",
			})...,
		)
		return errors.New("invalid order total")
	}

	log.Info("order created successfully",
		zap.Int64(logger.FieldDuration, duration.Milliseconds()),
		zap.String(logger.FieldStatus, "created"),
	)

	return nil
}

// GetOrders lista pedidos com paginação
func (s *OrderService) GetOrders(ctx context.Context, page, pageSize int) ([]Order, int, error) {
	log := logger.FromContext(ctx).With(
		zap.String(logger.FieldModule, "orders"),
		zap.String(logger.FieldFunction, "GetOrders"),
	)

	start := time.Now()

	log.Info("fetching orders",
		logger.WithPagination(page, pageSize, 0)...,
	)

	// Simular busca no banco de dados
	orders := []Order{
		{ID: "order-1", UserID: "user-123", Status: "pending", Items: 3, Total: 99.90},
		{ID: "order-2", UserID: "user-123", Status: "completed", Items: 1, Total: 49.90},
	}
	total := 2

	duration := time.Since(start)

	fields := []zap.Field{
		zap.Int(logger.FieldCount, len(orders)),
		zap.Int64(logger.FieldDuration, duration.Milliseconds()),
	}
	fields = append(fields, logger.WithPagination(page, pageSize, total)...)

	log.Info("orders fetched successfully", fields...)

	return orders, total, nil
}

// OrderHandler demonstra uso de logging em HTTP handlers
type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Logger já está no contexto graças ao middleware
	log := logger.FromContext(r.Context()).With(
		zap.String(logger.FieldModule, "handlers"),
		zap.String(logger.FieldFunction, "CreateOrderHandler"),
	)

	log.Info("processing create order request")

	// Simular parse do body
	order := &Order{
		ID:       "order-123",
		UserID:   logger.GetUserID(r.Context()), // Obtém user_id do contexto
		Status:   "pending",
		Items:    2,
		Total:    149.90,
		CreateAt: time.Now(),
	}

	// Chamar serviço
	if err := h.service.CreateOrder(r.Context(), order); err != nil {
		log.Error("failed to create order",
			zap.Error(err),
		)
		http.Error(w, "Failed to create order", http.StatusBadRequest)
		return
	}

	log.Info("order created successfully",
		zap.String("order_id", order.ID),
	)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": "%s", "status": "created"}`, order.ID)
}

func (h *OrderHandler) ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context()).With(
		zap.String(logger.FieldModule, "handlers"),
		zap.String(logger.FieldFunction, "ListOrdersHandler"),
	)

	// Simular parse de query params
	page := 1
	pageSize := 10

	orders, total, err := h.service.GetOrders(r.Context(), page, pageSize)
	if err != nil {
		log.Error("failed to fetch orders", zap.Error(err))
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	log.Info("orders listed successfully",
		zap.Int("count", len(orders)),
		zap.Int("total", total),
	)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"orders": %d, "total": %d}`, len(orders), total)
}

// Exemplo de função auxiliar que faz logging
func ProcessPayment(ctx context.Context, orderID string, amount float64) error {
	log := logger.FromContext(ctx).With(
		zap.String(logger.FieldModule, "payments"),
		zap.String(logger.FieldFunction, "ProcessPayment"),
		zap.String("order_id", orderID),
	)

	log.Info("processing payment",
		zap.Float64("amount", amount),
	)

	// Simular processamento
	time.Sleep(100 * time.Millisecond)

	// Dados sensíveis devem ser sanitizados
	creditCard := "1234-5678-9012-3456"
	log.Info("payment processed",
		zap.String("card", logger.SanitizeCreditCard(creditCard)),
		zap.String(logger.FieldStatus, "approved"),
	)

	return nil
}

func main() {
	// 1. Configurar logger
	log := logger.NewProduction()

	// Adicionar informações da aplicação
	log = logger.WithAppInfo(log, "order-service", "1.0.0")

	// Garantir que buffers sejam liberados ao sair
	defer log.Sync()

	log.Info("starting order service")

	// 2. Criar serviços e handlers
	orderService := &OrderService{}
	orderHandler := NewOrderHandler(orderService)

	// 3. Configurar rotas com middleware de logging
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", orderHandler.CreateOrderHandler)
	mux.HandleFunc("/orders/list", orderHandler.ListOrdersHandler)

	// 4. Aplicar middlewares
	var handler http.Handler = mux

	// Middleware de logging (adiciona request_id e loga requisições)
	handler = logger.LoggingMiddleware(log)(handler)

	// Middleware de recovery (captura panics)
	handler = logger.RecoveryMiddleware(log)(handler)

	// 5. Para health checks, pode-se pular o logging
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	healthMux.Handle("/", handler)

	// Usar SkipPathsMiddleware para pular /health
	finalHandler := logger.SkipPathsMiddleware(log, []string{"/health"})(healthMux)

	// 6. Exemplo de uso em contexto
	ctx := context.Background()
	ctx = logger.WithLogger(ctx, log)
	ctx = logger.WithUserID(ctx, "user-123")

	// Simular algumas operações
	order := &Order{
		ID:     "order-999",
		UserID: "user-123",
		Status: "pending",
		Items:  5,
		Total:  299.90,
	}

	if err := orderService.CreateOrder(ctx, order); err != nil {
		log.Error("failed to create test order", zap.Error(err))
	}

	if err := ProcessPayment(ctx, order.ID, order.Total); err != nil {
		log.Error("payment failed", zap.Error(err))
	}

	// 7. Iniciar servidor
	log.Info("server starting",
		zap.String("port", "8080"),
		zap.String("environment", "production"),
	)

	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatal("server failed to start", zap.Error(err))
	}
}

// Exemplos de output esperado no Grafana Loki:

/*
{
  "level": "info",
  "ts": "2025-10-24T10:30:45.123Z",
  "msg": "starting order creation",
  "environment": "production",
  "app_name": "order-service",
  "version": "1.0.0",
  "module": "orders",
  "function": "CreateOrder",
  "entity_type": "order",
  "entity_id": "order-123",
  "user_id": "user-123",
  "items_count": 2,
  "total": 149.90
}

{
  "level": "info",
  "ts": "2025-10-24T10:30:45.173Z",
  "msg": "order created successfully",
  "environment": "production",
  "app_name": "order-service",
  "version": "1.0.0",
  "module": "orders",
  "function": "CreateOrder",
  "entity_type": "order",
  "entity_id": "order-123",
  "duration_ms": 50,
  "status": "created"
}

{
  "level": "info",
  "ts": "2025-10-24T10:30:45.180Z",
  "msg": "request started",
  "environment": "production",
  "app_name": "order-service",
  "version": "1.0.0",
  "request_id": "req-abc123xyz",
  "http_method": "POST",
  "http_path": "/orders",
  "http_remote_addr": "192.168.1.100:54321",
  "http_host": "api.example.com",
  "http_user_agent": "Mozilla/5.0",
  "http_scheme": "https"
}

{
  "level": "info",
  "ts": "2025-10-24T10:30:45.250Z",
  "msg": "request completed",
  "environment": "production",
  "app_name": "order-service",
  "version": "1.0.0",
  "request_id": "req-abc123xyz",
  "http_method": "POST",
  "http_path": "/orders",
  "http_remote_addr": "192.168.1.100:54321",
  "http_status": 201,
  "duration_ms": 70,
  "http_response_size": 35
}
*/
