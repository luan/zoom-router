package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/luan/zoom-router/handlers"
)

var meetings handlers.Meetings

func main() {
	var meetingsPath string
	var port int

	flag.StringVar(&meetingsPath, "meetings", "meetings.json", "json file with a mappings of meeting names to zoom ids")
	flag.IntVar(&port, "port", 8080, "port to serve on")
	flag.Parse()

	if os.Getenv("PORT") != "" {
		var err error
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatal(err)
		}
	}

	meetingsFile, err := os.Open(meetingsPath)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(meetingsFile)
	err = dec.Decode(&meetings)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/", handlers.NewIndexHandler(meetings))
	r.Handle("/{meeting}", handlers.NewRedirectHandler(meetings))

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
