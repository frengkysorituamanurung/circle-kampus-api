package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/frengkysorituamanurung/circle-kampus-api/internal/handler"
	"github.com/frengkysorituamanurung/circle-kampus-api/internal/store"
)

func main() {
	// Sebaiknya ini datang dari file config atau environment variable
	dsn := "postgres://postgres:root123@localhost:3005/dev_circle_kampus?sslmode=disable"

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v\n", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal ping database: %v\n", err)
	}

	fmt.Println("🎉 Berhasil terhubung ke database PostgreSQL!")


	// Inisialisasi komponen
	
	userStore := store.NewUserStore(dbpool)
	userHandler := handler.NewUserHandler(userStore)

	// Inisialisasi Gin Router
	router := gin.Default()

	// Grouping routes by version
	v1 := router.Group("/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			// Rute login akan kita tambahkan di sini nanti
		}
	}
	
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	

	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	router.Run(":8080") 
}