// Package tee provides a multiplexing io.Writer that duplicates every
// write to a set of underlying destinations.
//
// It is useful when cronlog output should be emitted to more than one
// sink simultaneously — for example, writing structured JSON to a file
// while also streaming human-readable text to stdout.
//
// Basic usage:
//
//	file, _ := os.Create("job.log")
//	w, err := tee.New(os.Stdout, file)
//	if err != nil {
//		log.Fatal(err)
//	}
//	logger := logger.New("my-job", w)
//
// Write errors from individual destinations do not abort the fan-out;
// all destinations are always attempted and any errors are joined and
// returned to the caller.
package tee
