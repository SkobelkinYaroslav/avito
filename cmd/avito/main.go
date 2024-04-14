package main

import (
	"avito/internal/handler"
	"avito/internal/repository"
	"avito/internal/service"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalf("Error while closing server %q", err)
		}
	}(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	authRepo := repository.NewAuthRepository(db)
	bannerRepo := repository.NewBannerRepository(db, rdb)

	authService := service.NewAuthService(authRepo)
	bannerService := service.NewBannerService(bannerRepo)

	router := handler.NewHandler(authService, bannerService)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error running server: %q", err)
	}
}
