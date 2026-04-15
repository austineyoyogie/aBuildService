package server

import (
	"aIBuildService/aPI/config"
	"aIBuildService/aPI/config/database"
	"aIBuildService/aPI/resource"
	"aIBuildService/aPI/service/implementation"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

var (
	CD = database.CDriver()
	LC = config.LoadConfig()
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Start() error {
	router := http.NewServeMux()
	router.Handle("/v3/", http.StripPrefix("/v3", router))

	userServiceImpl := implementation.NewUserServiceImpl(CD)
	userResource := resource.NewUserResource(userServiceImpl, LC.JWT.KEY)
	userResource.UserResourceHandlerRoutes(router)

	logger := CustomRouterLogger("RouterLoggerLogs.txt")
	server := http.Server{
		Addr:    s.addr,
		Handler: loadCors(LoggerMiddleware(logger, router)),
	}
	log.Printf("\nStarting [SERVER] on Port %s", s.addr)
	return server.ListenAndServe()
}

func loadCors(r http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:3000"},
		AllowedHeaders:   []string{"Origin", "Access-Control-Allow-Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Location", "Entity"},
		ExposedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "true"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
	return corsMiddleware.Handler(r)
}

func LoggerMiddleware(logger *log.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Request received: Method %s, Path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func CustomRouterLogger(filePath string) *log.Logger {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	routerLogger := log.New(file, "Router Logger: ", log.LstdFlags|log.Lshortfile)
	return routerLogger
}
