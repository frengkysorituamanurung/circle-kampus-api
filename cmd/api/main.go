package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/frengkysorituamanurung/circle-kampus-api/configs"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/handler"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/store"
)

func main() {

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Gagal memuat file .env: %v\n", err)
	}

	dsn := os.Getenv("DATABASE_URL")

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v\n", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal ping database: %v\n", err)
	}

	fmt.Println("🎉 Berhasil terhubung ke database PostgreSQL!")


	
	userStore := store.NewUserStore(dbpool)
	userHandler := handler.NewUserHandler(userStore)

	router := gin.Default()

	configs.SetupCORS(router)

	v1 := router.Group("/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}
		protectedRoutes := v1.Group("/", handler.AuthMiddleware())
		{
			protectedRoutes.GET("/users/me", userHandler.GetProfile)
		}
	}
	

	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	router.Run(":8080") 
}
