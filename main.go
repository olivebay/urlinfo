package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"github.com/olivebay/urlinfo/api/models"
	"github.com/olivebay/urlinfo/api/handlers"
)

// Parse the configuration file 'config.toml', and establish a connection to DB

func main() {

	// todo maybe add a common config file
	var (
		bindAddress = flag.String("addr", ":9090", "endpoint address")
	)

	l := log.New(os.Stdout, "urlinfo-api ", log.LstdFlags)

	// Create and store a Mongo session for every requests.
	session := models.NewSession()
	viper.Set("MONGO_SESSION", session)
	defer session.Close()

	// create a new serve mux and register the handlers
	r := mux.NewRouter().SkipClean(true).UseEncodedPath()
	r.Use(MongoHandler) // 
	r.Use(mux.CORSMethodMiddleware(r)) // CORS middleware

	// handlers for API
	r.HandleFunc("/healthz", handlers.StatusHandler)
	api := r.PathPrefix(`/urlinfo/1/`).Subrouter()
	api.HandleFunc(`/{url:.+}`, handlers.GetURL).Methods("Get", "Head")

	// TODO update database

	// create a new server
	srv := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println("Starting server on port", *bindAddress)

		err := srv.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received
	sig := <-c
	log.Println("Got signal:", sig)

	// Create a deadline to wait for
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(ctx)
}


// MongoHandler insert Mgo.session in context and serve the request.
func MongoHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbSession := viper.Get("MONGO_SESSION").(models.Session).Copy()
		r = r.WithContext(
			context.WithValue(r.Context(), "db", dbSession))
		next.ServeHTTP(w, r)
		dbSession.Close()
	})
}