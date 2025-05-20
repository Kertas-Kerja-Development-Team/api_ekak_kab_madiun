package main

import (
	"ekak_kabupaten_madiun/helper"
	"ekak_kabupaten_madiun/middleware"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	host := os.Getenv("host")
	port := os.Getenv("port")
	addr := fmt.Sprintf("%s:%s", host, port)

	cors := helper.NewCORSMiddleware()

	if addr == ":" {
		addr = "localhost:8080"
	}

	return &http.Server{
		Addr:    addr,
		Handler: cors.Handler(authMiddleware),
	}
}

func main() {
	runSeeder := flag.Bool("seed", false, "Jalankan database seeder")
	flag.Parse()
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Cek flag seeder
	if *runSeeder {
		log.Println("Menjalankan database seeder...")
		seeder := InitializeSeeder()
		seeder.SeedAll()
		log.Println("Seeder selesai dijalankan")
		return
	}
	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL") // e.g. https://keycloak.local/realms/myrealm/protocol/openid-connect/certs
	erru := middleware.InitJWKS(jwksURL)
	if erru != nil {
		log.Fatalf("Failed to initialize JWKS: %v", err)
	}

	// Initialize dan jalankan server
	server := InitializeServer()
	log.Printf("Server berjalan di %s", server.Addr)
	err = server.ListenAndServe()
	helper.PanicIfError(err)
}
