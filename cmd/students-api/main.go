package main

import (
	"context"
	
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SudipSarkar1193/students-API-Go/internal/config"

	"net/http"
)

func main() {
	// load config

	cfg :=config.MustLoad()

	//database setup

	//setup router
     
	router := http.NewServeMux()

	router.HandleFunc("GET /",func(rw http.ResponseWriter, r *http.Request){ //⭐⭐ "GET /" ->The space after the GET is mandatory 
		rw.Write([]byte("Welcome to students api"))
	})

	//setup server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at",slog.String("PORT",cfg.Addr))


	done :=make(chan os.Signal,1)

	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func ()  {
		err := server.ListenAndServe() 
		if err!=nil{
			log.Fatal("Failed to start server")
		}
	}()

	<-done 

	//Logic to stop server 

	slog.Info("Shutting down the server...")

	ctx,cancel :=context.WithTimeout(context.Background(),5*time.Second)

	defer cancel()


	err := server.Shutdown(ctx)
	if err!=nil {
		slog.Error("Failed to shutdown server",slog.String("Error : ",err.Error()))
	}

	slog.Info("Server shut down successfully")
}
