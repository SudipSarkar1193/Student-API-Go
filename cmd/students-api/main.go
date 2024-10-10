package main

import (
	"context"
	"fmt"

	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SudipSarkar1193/students-API-Go/internal/config"
	"github.com/SudipSarkar1193/students-API-Go/internal/storage/mySql_Db"
	"github.com/joho/godotenv"

	"github.com/SudipSarkar1193/students-API-Go/internal/http/handlers"
)

func main() {

	//

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// load config

	cfg := config.MustLoad()

	//database setup

	storage, err := mySql_Db.NewOnline(cfg)
	if err != nil {
		log.Fatal("DB connecton ERROR : ", err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env))

	fmt.Println("storage : ", storage)

	//setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students/new", student.New(storage))
	router.HandleFunc("GET /api/students/all", student.GetAllStudents(storage))
	router.HandleFunc("POST /api/students/login", student.Login(storage))
	//router.HandleFunc("POST /api/student", student.AddIsSafeMiddleware(student.GetStudentsByIdOrEmail(storage)))

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at", slog.String("PORT", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	//Logic to stop server

	slog.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("Error : ", err.Error()))
	}

	slog.Info("Server shut down successfully")
}
