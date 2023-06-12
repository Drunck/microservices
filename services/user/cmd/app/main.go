package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"microservices-go/pkg/hash"
	"microservices-go/pkg/store/postgres"
	"microservices-go/pkg/token"
	"microservices-go/services/user/internal/delivery/http"
	"microservices-go/services/user/internal/repository"
	"microservices-go/services/user/internal/service"
	"os"
)

func main() {
	dbConnCfg := postgres.ConnConfig{}
	httpServerCfg := http.ServerConfig{}

	flag.IntVar(&httpServerCfg.Port, "http-port", 4000, "HTTP server port")
	flag.StringVar(&httpServerCfg.ReadTimeout, "http-read-timeout", "10s", "HTTP read timeout")
	flag.StringVar(&httpServerCfg.WriteTimeout, "http-write-timeout", "30s", "HTTP write timeout")
	flag.StringVar(&httpServerCfg.IdleTimeout, "http-idle-timeout", "1m", "HTTP idle timeout")

	flag.IntVar(&dbConnCfg.Port, "pg-port", 5432, "Postgres port")
	flag.StringVar(&dbConnCfg.Host, "pg-host", "localhost", "Postgres host")
	flag.StringVar(&dbConnCfg.User, "pg-user", os.Getenv("PG_USER"), "Postgres user")
	flag.StringVar(&dbConnCfg.Password, "pg-password", os.Getenv("PG_PASSWORD"), "Postgres password")
	flag.StringVar(&dbConnCfg.DbName, "pg-db-name", os.Getenv("PG_DB_NAME"), "Postgres DB name")
	flag.IntVar(&dbConnCfg.MaxOpenConns, "pg-max-open-conns", 15, "Postgres max open connections")
	flag.StringVar(&dbConnCfg.MaxIdleTime, "pg-max-idle-time", "15m", "Postgres max connection idle time")
	flag.Parse()

	db, err := postgres.OpenDB(dbConnCfg)
	if err != nil {
		log.Err(err)
		os.Exit(1)
	}
	defer db.Pool.Close()

	log.Print("database connection pool established")

	passwordsHashCost := hash.NewBCryptHasher(12)
	tokenManager, err := token.NewManager(os.Getenv("TOKEN_KEY"))
	userRepository := repository.NewUserRepo(db.Pool)
	userService := service.New(userRepository, passwordsHashCost, tokenManager)

	httpServer := http.NewHttpServer(http.NewRouter(userService).GetRoutes(), httpServerCfg)

	err = httpServer.Serve()
	if err != nil {
		log.Err(err).Msg("Failed to start HTTP server")
	}

}
