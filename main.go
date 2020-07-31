package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func main(){
	fmt.Println("hogehgoe")

	// Create your application using your preferred app name, license key, and
	// any other configuration options.
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("isucon_newrelic_sample"),
		newrelic.ConfigLicense(os.Getenv("KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	// Now you can use the Application to collect data!  Create transactions
	// to time inbound requests or background tasks. You can start and stop
	// transactions directly using Application.StartTransaction and
	// Transaction.End.
	func() {
		txn := app.StartTransaction("myTask")
		defer txn.End()

		// Do some work
		time.Sleep(time.Second)
	}()

	// WrapHandler and WrapHandleFunc make it easy to instrument inbound
	// web requests handled by the http standard library without calling
	// Application.StartTransaction.  Popular framework instrumentation
	// packages exist in the v3/integrations directory.
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/hoge", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "this is the index page")
	}))

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/hello", func(w http.ResponseWriter, req *http.Request) {
		// WrapHandler and WrapHandleFunc add the transaction to the
		// inbound request's context.  Access the transaction using
		// FromContext to add attributes, create segments, and notice.
		// errors.
		txn := newrelic.FromContext(req.Context())

		func() {
			// Segments help you understand where the time in your
			// transaction is being spent.  You can use them to time
			// functions or arbitrary blocks of code.
			defer txn.StartSegment("helperFunction").End()
		}()

		io.WriteString(w, "hello world")
	}))

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/slow", func(w http.ResponseWriter, req *http.Request) {
		// Do some work
		time.Sleep(time.Second * 5)

		io.WriteString(w, "this is slow page")
	}))

	http.ListenAndServe(":8000", nil)
}