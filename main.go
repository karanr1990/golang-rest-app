package main

import (
	"context"
	"github.ibm.com/Quest-CIO/go-micro-app/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main()  {
	l := log.New(os.Stdout,"product-api",log.LstdFlags)

	hh :=handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()

	sm.Handle("/",hh)
	sm.Handle("/goodbye",gh)
	sm.Handle("/products",ph)


	s := &http.Server{
		Addr: ":4000",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}

	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan,os.Interrupt)
	signal.Notify(sigchan,os.Kill)

	sig := <- sigchan

	l.Println("recieved terminate...",sig)

	tc, _ := context.WithTimeout(context.Background(),30*time.Second)

	s.Shutdown(tc)
}

