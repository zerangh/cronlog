// Package circuitbreaker implements a thread-safe circuit breaker pattern
// suitable for protecting external calls — such as webhook notifications —
// made during cron job execution.
//
// # States
//
// The circuit breaker moves through three states:
//
//   - Closed: all calls are allowed through; failures are counted.
//   - Open: calls are rejected immediately with ErrOpen until the reset
//     timeout elapses.
//   - HalfOpen: one probe call is permitted; a success closes the circuit
//     while a failure reopens it.
//
// # Usage
//
//	breaker, err := circuitbreaker.New(5, 30*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := breaker.Allow(); err != nil {
//		// circuit is open — skip the call
//		return err
//	}
//
//	if err := doExternalCall(); err != nil {
//		breaker.RecordFailure()
//		return err
//	}
//	breaker.RecordSuccess()
package circuitbreaker
