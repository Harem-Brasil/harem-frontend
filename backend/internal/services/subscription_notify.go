package services

import (
	"context"
	"database/sql"

	"github.com/harem-brasil/backend/internal/realtime"
)

func (s *Services) fanoutSubscriptionUpdated(ctx context.Context, subscriberUserID string, creatorUserID *string, subscriptionID, status string) {
	pub := s.RealtimePublisher
	if pub == nil {
		pub = realtime.NoopPublisher{}
	}
	payload := realtime.SubscriptionUpdatedPayload{
		SubscriptionID: subscriptionID,
		Status:         status,
	}
	recipients := []string{subscriberUserID}
	if creatorUserID != nil && *creatorUserID != "" && *creatorUserID != subscriberUserID {
		recipients = append(recipients, *creatorUserID)
	}
	if err := pub.PublishSubscriptionUpdated(ctx, recipients, payload); err != nil && s.Logger != nil {
		s.Logger.Warn("realtime subscription.updated publish failed", "error", err.Error())
	}
}

func nullableString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	v := ns.String
	return &v
}
