package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"github.com/pkg/errors"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := reflect.ValueOf(http.DefaultServeMux).Elem().FieldByName("m")

		w.WriteHeader(http.StatusOK)
		for _, path := range m.MapKeys() {
			fmt.Fprintln(w, path)
		}
	})

	http.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("omg")
	})

	http.HandleFunc("/stdout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mw := io.MultiWriter(w, os.Stdout)
		fmt.Fprintf(mw, "okok")
	})

	http.HandleFunc("/stderr", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "owf")
	})

	http.HandleFunc("/debug.Stack", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "%s", debug.Stack())
	})

	http.HandleFunc("/github.com/pkg/errors", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("uwaaaa")

		w.WriteHeader(http.StatusInternalServerError)
		mw := io.MultiWriter(w, os.Stderr)
		fmt.Fprintf(mw, "%+v", err)
	})

	http.HandleFunc("/errorreporting", func(w http.ResponseWriter, r *http.Request) {
		projectID := "pokutuna-dev"
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
		fmt.Fprintf(w, "%+v", appErr)
	})

	http.ListenAndServe(":8080", nil)
}
