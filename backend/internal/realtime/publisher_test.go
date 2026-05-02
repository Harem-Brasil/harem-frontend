package realtime

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNoopPublisher_PublishSubscriptionUpdated(t *testing.T) {
	var p NoopPublisher
	err := p.PublishSubscriptionUpdated(t.Context(), []string{"u1"}, SubscriptionUpdatedPayload{
		SubscriptionID: "550e8400-e29b-41d4-a716-446655440000",
		Status:         "active",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPayloadPublicJSON_noSensitiveKeys(t *testing.T) {
	b, err := PayloadPublicJSON(SubscriptionUpdatedPayload{
		SubscriptionID: "550e8400-e29b-41d4-a716-446655440001",
		Status:         "canceled",
	})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	forbidden := []string{"price", "amount", "currency", "customer", "email", "payment", "card"}
	lower := strings.ToLower(string(b))
	for _, w := range forbidden {
		if strings.Contains(lower, w) {
			t.Fatalf("serialized envelope should not contain %q", w)
		}
	}
	if m["v"] == nil || m["type"] == nil || m["payload"] == nil {
		t.Fatalf("expected v, type, payload keys, got %v", m)
	}
}
