package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/martbul/internal/db"
	"github.com/martbul/pkg/logger"
)

func main() {
	logger.Init()
	log := logger.GetLogger()

	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file")
	}

	dbConn, err := db.ConnectDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	surveMux := mux.NewRouter()

	surveMux.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	//TODO: Move to config dir
	//	duelMM := duelMatchmaker.NewDuelMatchmaker(2)
	///	streakGameMM := streakGameMatchmaker.NewStreakGameMatchmaker()
	//	pinPointSPGameMM := pinPointSPGameMatchmaker.NewPinPointSPGameMatchmaker()
	//	api := surveMux.PathPrefix("").Subrouter()
	//	user.RegisterUserRoutes(api, dbConn)
	//	game.RegisterGameRoutes(api, duelMM, streakGameMM, pinPointSPGameMM, dbConn)
	//	stats.RegisterStatsRoutes(api, dbConn)

	cors := gohandlers.CORS(
		gohandlers.AllowedOrigins([]string{"*"}),
		gohandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gohandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		gohandlers.AllowCredentials(),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      cors(surveMux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// Set up signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Block until a signal is received
	sig := <-sigChan
	log.Info("Received shutdown signal:", sig)

	// Create a context with timeout for graceful shutdown
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(timeoutContext); err != nil {
		log.Error("Server shutdown failed", "error", err)
	} else {
		log.Info("Server shut down gracefully")
	}
}
