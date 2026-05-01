package services

import (
	"log/slog"

	"github.com/harem-brasil/backend/internal/realtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Dependencies agrega infraestrutura injetável nas regras de negócio.
type Dependencies struct {
	DB          *pgxpool.Pool
	Redis       *redis.Client
	JWTSecret   []byte
	Logger      *slog.Logger
	MaxFileSize int64

	StripeWebhookSecret      string
	PagSeguroWebhookSecret   string
	MercadoPagoWebhookSecret string
	// InternalBillingSecret protege callbacks internos (fila/worker → marcar pedido pago). Vazio em dev/test pode ser aceite só em ValidateInternalBillingSecret.
	InternalBillingSecret string
	// AppEnv replica ENV (ex.: development, test, production); usado quando o segredo do webhook está vazio.
	AppEnv string
	// PlatformCommissionBasisPoints comissão da plataforma sobre pedidos pagos do catálogo (0–10000). Predefinido em composition root (ex.: 1500 = 15%).
	PlatformCommissionBasisPoints int
	// RealtimePublisher hub WS (HB-EPIC-04); nil usa noop até existir implementação.
	RealtimePublisher realtime.Publisher
}
