// Package ratelimit provides a token-bucket rate limiter intended for use
// with webhook notifications in cronlog.
//
// # Overview
//
// When a cron job emits many errors in rapid succession it can flood an
// alerting endpoint. The Limiter type controls how many notification
// attempts are allowed within a given time window.
//
// # Usage
//
//	l, err := ratelimit.New(ratelimit.Config{
//		Max:       5,   // burst of up to 5 notifications
//		PerSecond: 0.5, // replenish one token every 2 seconds
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = l.Do(func() error {
//		return notifier.Notify(payload)
//	})
//	if errors.Is(err, ratelimit.ErrRateLimited) {
//		// notification was suppressed
//	}
package ratelimit
