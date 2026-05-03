package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/client"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/repository/postgres"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/server"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/service"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("warning: assuming default configuration: .env unreadable: %v\n", err)
	}

	cfg := config.MustLoad()

	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	cityRepo := postgres.NewCityRepository(db)
	weatherRepo := postgres.NewWeatherRepository(db)

	userService := service.NewUserService(userRepo)
	cityService := service.NewCityService(cityRepo)
	weatherService := service.NewWeatherService(weatherRepo)

	osmClient := client.NewOsmClient(cfg.Api)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherClient := client.NewWeatherClient(httpClient)

	jwtManager := auth.NewJWTManager([]byte(cfg.App.JwtSecret))

	handler := server.NewHandler(userService, cityService, weatherService,
		osmClient, weatherClient, jwtManager)
	router := server.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	go func() {
		log.Printf("server started on :%s", cfg.App.Port)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("shutdown server: %v", err)
	}

	log.Println("server stopped")
}
