package middleware

import (
	"context"
	"time"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"go.uber.org/zap"
)

type middlewareWrapper func(next tg.Invoker) telegram.InvokeFunc

func (m middlewareWrapper) Handle(next tg.Invoker) telegram.InvokeFunc {
	return m(next)
}

func FloodWait(logger *zap.Logger) telegram.Middleware {
	return middlewareWrapper(func(next tg.Invoker) telegram.InvokeFunc {
		return telegram.InvokeFunc(func(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
			for {
				err := next.Invoke(ctx, input, output)
				if err == nil {
					return nil
				}

				d, ok := tgerr.AsFloodWait(err)
				if !ok {
					return err
				}

				if d > 5*time.Minute {
					logger.Warn("FloodWait too long, aborting", zap.Duration("duration", d))
					return err
				}

				logger.Info("FloodWait detected, sleeping...", zap.Duration("duration", d))

				timer := time.NewTimer(d + 1*time.Second)
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
					continue
				}
			}
		})
	})
}
