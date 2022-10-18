//Filename: cmd/api/server.go

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {

	//http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//shutdown func should return its errors to its channel

	shutdownError := make(chan error)

	//start a background go routine

	go func() {
		//create a quite/exit channel which carries os.Signal values
		quit := make(chan os.Signal, 1)
		//listen for sigint and sigterm
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		//block until a signal is received
		s := <-quit

		//log a message

		app.logger.PrintInfo("Shutting down server", map[string]string{
			"signal": s.String(),
		})

		//create a context  with a 20 sec timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		//call the shutdown function
		shutdownError <- srv.Shutdown(ctx)

	}()

	//Start Server
	app.logger.PrintInfo("Starting  server on ", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	// check if the the shutdown process has initiated
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	//Block for notification from the shutdown function

	err = <-shutdownError
	if err != nil {
		return err
	}

	//Graceful shutdown was successful
	app.logger.PrintInfo("stopped server ", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
