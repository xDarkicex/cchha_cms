package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/xDarkicex/todos"

	"github.com/NYTimes/gziphandler"
	"github.com/xDarkicex/cchha_server_new/helpers"
	"github.com/xDarkicex/cchha_server_new/server"
	"golang.org/x/crypto/acme/autocert"
)

var isProduction *bool
var todos elephant.Elephant

func init() {
	todos = elephant.NewElephant(elephant.SetName("todos"), elephant.WithPersistance(true))
	todos.SetMemory("Todo", "Build logger middleware.")
	todos.SetMemory("Todo", "Figure out why emit emits so many times.")
	todos.SetMemory("Todo", "Fix Admin signup page, Javascript is broken for lock")
	todos.SetMemory("Todo", "Design signup page too look like signin")
	isProduction = flag.Bool("Production", false, "")
	flag.Parse()

	todos.SaveMemories()
}

func main() {

	// cert := "/etc/letsencrypt/live/cchha.com/fullchain.pem"
	// key := "/etc/letsencrypt/live/cchha.com/privkey.pem"

	if false {
		domains := []string{"compassionatecare.com", "www.compassionatecare.com", "cchha.com", "www.cchha.com"}
		fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU()))
		handler := server.NewRouter()
		gzipped := gziphandler.GzipHandler(handler)
		srv := &http.Server{
			Addr:         ":443",
			Handler:      gzipped,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		go func() {
			helpers.HandleError((http.ListenAndServe(":80", http.HandlerFunc(helpers.HTTPS))))
		}()
		// Start the server
		helpers.HandleError(srv.Serve(autocert.NewListener(domains...)))
		// Wait for an interrupt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		// Attempt a graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println(helpers.Red("server shutdown"))
		helpers.HandleError(srv.Shutdown(ctx))
	}
	// development server
	fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU()))
	handler := server.NewRouter()
	// gzipped := gziphandler.GzipHandler(handler)
	srv := &http.Server{
		Addr:         ":3001",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	helpers.HandleError((srv.ListenAndServe()))
}
