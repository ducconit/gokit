package otp

import (
	"context"
	"fmt"
)

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

type Router struct {
	Senders map[Channel]Sender
}

func (r Router) Send(ctx context.Context, msg Message) error {
	s := r.Senders[msg.Channel]
	if s == nil {
		return fmt.Errorf("otp: missing sender for channel %s", msg.Channel)
	}
	return s.Send(ctx, msg)
}
