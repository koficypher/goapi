package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// http request handlers
	r.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, "Hello Gophers")
	})

	// constructing a typical server using the http server struct
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
		Handler:      r,
	}

	fmt.Println("Goapi server running .....")

	// make an error channel of size 2 to collect all errors that may arise from starting the server
	errs := make(chan error, 2)

	// goroutine to start the server concurrently
	go func() {
		// start server from the previously constructed server
		if err := s.ListenAndServe(); err != nil {
			errs <- err
		}
	}()

	// goroutine that creates a signal channel that sends signal to the errors channel when something happens
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)

		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Shutting down API server with error: %s", <-errs)

	// start http serve directly
	//	http.ListenAndServe(":8080", nil)

	// perform a graceful shutdown of the server by giving a 5 seconds time limit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// defer the cancel function
	defer cancel()

	// shutdown the server gracefully
	s.Shutdown(ctx)
}
