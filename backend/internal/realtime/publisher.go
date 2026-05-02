package realtime

import (
	"context"
	"encoding/json"
)

// Documentado em API_IMPLEMENTACAO.md secção 8 (eventos servidor → cliente).
const ServerEventSubscriptionUpdated = "subscription.updated"

// SubscriptionUpdatedPayload é o payload público WS — apenas IDs e estado, sem valores monetários nem PII.
type SubscriptionUpdatedPayload struct {
	SubscriptionID string `json:"subscription_id"`
	Status         string `json:"status"`
}

// EnvelopeWS envelope mínimo alinhado à secção 8 (v=1, type, payload).
type EnvelopeWS struct {
	V       int                        `json:"v"`
	Type    string                     `json:"type"`
	Payload SubscriptionUpdatedPayload `json:"payload"`
}

// Publisher publica eventos para utilizadores autenticados no hub (fan-out por user_id).
// Implementação real: HB-EPIC-04 (hub WS + Redis entre instâncias).
type Publisher interface {
	PublishSubscriptionUpdated(ctx context.Context, recipientUserIDs []string, payload SubscriptionUpdatedPayload) error
}

// NoopPublisher satisfaz M7 sem servidor WS — não envia tráfego.
type NoopPublisher struct{}

func (NoopPublisher) PublishSubscriptionUpdated(ctx context.Context, recipientUserIDs []string, payload SubscriptionUpdatedPayload) error {
	_ = ctx
	_ = recipientUserIDs
	_ = payload
	return nil
}

// PayloadPublicJSON valida que o JSON serializado não inclui campos extra não previstos (testes).
func PayloadPublicJSON(p SubscriptionUpdatedPayload) ([]byte, error) {
	env := EnvelopeWS{V: 1, Type: ServerEventSubscriptionUpdated, Payload: p}
	return json.Marshal(env)
}
