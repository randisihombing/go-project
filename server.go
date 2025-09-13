package main

import (
	"crypto/tls"
	"fmt"
	mw "gocourse/internal/api/middlewares"
	"gocourse/internal/repository/sqlconnect"
	"gocourse/internal/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		return
	}

	_, err = sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println("Error------:", err)
		return
	}
	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))

	// secureMux := utils.ApplyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)

	router := router.Router()
	secureMux := mw.SecurityHeaders(router)

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running in port: ", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
}
