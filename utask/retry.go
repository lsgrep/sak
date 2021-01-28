package utask

import (
	"context"
	"github.com/lsgrep/sak/ulog"
	"github.com/lsgrep/sak/utime"
	"time"
)

// Returns nil on the first time `f()` returns nil, even if by that time `ctx` has
// been cancelled. Otherwise returns the last error returned by `f()`. If `f()` has
// never got a chance to run, returns `ctx.Err()`.
func TryUntilNoErr(ctx context.Context, funcName string, f func() error) (retErr error) {
	logCooldown := utime.NewUnreadyCooldown(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			if retErr == nil {
				return ctx.Err()
			}
			return retErr

		default:
			retErr = f()
			if retErr == nil {
				return nil
			}

			if logCooldown.Ready() {
				ulog.Errorf("%s failed with err: %v", funcName, retErr)
			}
			utime.SleepMs(200) // avoid busy retry
		}
	}
}
