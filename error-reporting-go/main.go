package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"github.com/pkg/errors"
)

var projectID = "pokutuna-playground"

func main() {
	// list endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m")

		w.WriteHeader(http.StatusOK)
		for _, path := range m.MapKeys() {
			fmt.Fprintln(w, path)
		}
	})

	// collected on Standard, the log entry is formatted
	// NOT collected on Flexible
	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("omg")
	})

	// NOT collected
	http.HandleFunc("/stdout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mw := io.MultiWriter(w, os.Stdout)
		fmt.Fprintf(mw, "okok\n")
	})

	// NOT collected
	http.HandleFunc("/stderr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "owf\n")
	})

	// NOT collected
	http.HandleFunc("/debug.Stack", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "%s", debug.Stack())
	})

	// NOT collected
	http.HandleFunc("/github.com/pkg/errors", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("uwaaaa")

		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "%+v", err)
	})

	// NOT collected
	http.HandleFunc("/runtime.Stack", func(w http.ResponseWriter, r *http.Request) {
		// it's same format with debug.Stack
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		fmt.Fprintf(mw, "%s\n", buf[:n])
	})

	// collected, stacktrace works on console
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		prop := r.URL.Query().Get("prop")
		if prop == "" {
			prop = "message"
		}

		payload := map[string]interface{}{
			prop: string(debug.Stack()),
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	// collected, stacktrace works on console
	http.HandleFunc("/json/with-type", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		payload := map[string]interface{}{
			"@type":   "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent",
			"message": string(debug.Stack()),
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	// collected, stacktrace works on console
	http.HandleFunc("/json/with-logName", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		payload := map[string]interface{}{
			"logName": fmt.Sprintf("projects/%s/clouderrorreporting.googleapis.com%%2Freported_errors", projectID),
			"message": string(debug.Stack()),
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	// collected, stacktrace works on console
	http.HandleFunc("/json/with-pkg-errors", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		err := errors.New("made by pkg/errors")

		payload := map[string]interface{}{
			"@type":   "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent",
			"message": fmt.Sprintf("%+v", err),
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	// collected, stacktrace is not available
	http.HandleFunc("/json/with-type/not-stacktrace", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		payload := map[string]interface{}{
			"@type":   "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent",
			"message": "hello with @type field",
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	// NOT collected
	// the logName field is not lited up to metadata on LogEntry.
	// it's 'projects/${projectID}/logs/stderr'
	http.HandleFunc("/json/with-logName/not-stacktrace", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)

		payload := map[string]interface{}{
			"logName": fmt.Sprintf("projects/%s/clouderrorreporting.googleapis.com%%2Freported_errors", projectID),
			"message": "hello with logName field",
		}
		output, _ := json.Marshal(payload)
		fmt.Fprintf(mw, "%s\n", output)
	})

	http.HandleFunc("/errorreporting", func(w http.ResponseWriter, r *http.Request) {
		projectID := "pokutuna-playground"
		errorClient, err := errorreporting.NewClient(context.Background(), projectID, errorreporting.Config{
			ServiceName: "error-reporting-go",
		})
		if err != nil {
			fmt.Fprintf(os.Stdout, "%+v", err)
		}
		defer errorClient.Close()

		appErr := errors.New("ohyooo")
		errorClient.Report((errorreporting.Entry{
			Error: appErr,
		}))
		fmt.Fprintf(w, "%+v\n", appErr)
	})

	http.ListenAndServe(":8080", nil)
}
