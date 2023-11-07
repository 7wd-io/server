package dispatcher

import (
	"7wd.io/domain"
	"context"
	"log/slog"
)

func New() *D {
	return &D{
		handlers: map[domain.EventId][]H{},
	}
}

type D struct {
	handlers map[domain.EventId][]H
}

func (dst *D) On(event domain.EventId, h ...H) *D {
	dst.handlers[event] = append(dst.handlers[event], h...)

	return dst
}

func (dst *D) Dispatch(ctx context.Context, event domain.EventId, payload interface{}) {
	for k, h := range dst.handlers[event] {
		if err := h(ctx, payload); err != nil {
			slog.Error(
				"dispatcher: "+err.Error(),
				slog.Int("event", int(event)),
				slog.Int("handler index", k),
			)
		}
	}
}

type H func(ctx context.Context, payload interface{}) error
