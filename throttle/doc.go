// Package throttle provides a concurrency-limiting primitive for use in
// cronlog pipelines.
//
// A Throttle restricts how many goroutines may simultaneously perform a
// guarded operation. This is useful when processing cron job output that
// may fan out to multiple downstream handlers (webhook delivery, log
// rotation, buffer writes) and you want to cap resource usage.
//
// Basic usage:
//
//	th, err := throttle.New(4)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := th.Acquire(ctx); err != nil {
//		// context cancelled — skip or return error
//		return err
//	}
//	defer th.Release()
//
//	// ... perform guarded work ...
//
// Throttle is safe for concurrent use by multiple goroutines.
package throttle
