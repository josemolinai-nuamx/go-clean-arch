package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	mysqlRepo "github.com/josemolinai-nuamx/go-clean-arch/internal/repository/mysql"

	"github.com/joho/godotenv"
	"github.com/josemolinai-nuamx/go-clean-arch/article"
	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest"
	"github.com/josemolinai-nuamx/go-clean-arch/internal/rest/middleware"
)

const (
	defaultTimeout         = 15
	defaultAddress         = ":9090"
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = 300
	defaultConnMaxIdleTime = 120
)

func envInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		log.Printf("invalid %s value (%q), using default %d", key, raw, fallback)
		return fallback
	}

	return value
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}

	maxOpenConns := envInt("DB_MAX_OPEN_CONNS", defaultMaxOpenConns)
	maxIdleConns := envInt("DB_MAX_IDLE_CONNS", defaultMaxIdleConns)
	connMaxLifetimeSec := envInt("DB_CONN_MAX_LIFETIME_SEC", defaultConnMaxLifetime)
	connMaxIdleTimeSec := envInt("DB_CONN_MAX_IDLE_TIME_SEC", defaultConnMaxIdleTime)

	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)
	dbConn.SetConnMaxLifetime(time.Duration(connMaxLifetimeSec) * time.Second)
	dbConn.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeSec) * time.Second)

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()
	// prepare echo

	e := echo.New()
	e.Use(middleware.CORS)
	timeout := envInt("CONTEXT_TIMEOUT", defaultTimeout)
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repository
	authorRepo := mysqlRepo.NewAuthorRepository(dbConn)
	articleRepo := mysqlRepo.NewArticleRepository(dbConn)

	// Build service Layer
	svc := article.NewService(articleRepo, authorRepo)
	rest.NewArticleHandler(e, svc)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address)) //nolint
}
