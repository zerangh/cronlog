// Package dispatcher provides a multi-destination log entry router for cronlog.
//
// A Dispatcher holds one or more named Routes. Each Route pairs an io.Writer
// with a minimum log level. When Dispatch is called, the entry is forwarded
// only to routes whose minimum level is satisfied by the entry's level.
//
// Level ordering (lowest to highest): debug < info < warn < error.
//
// Example usage:
//
//	errFile, _ := os.Create("errors.log")
//	allFile, _ := os.Create("all.log")
//
//	d, err := dispatcher.New(
//		dispatcher.NewRoute("all",    allFile,  "debug"),
//		dispatcher.NewRoute("errors", errFile,  "error"),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	d.Dispatch(dispatcher.Entry{Level: "info",  Message: "job started"})
//	d.Dispatch(dispatcher.Entry{Level: "error", Message: "job failed"})
package dispatcher
