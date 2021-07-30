package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"

	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/khihadysucahyo/go-clean-arch-boilerplate/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/khihadysucahyo/go-clean-arch-boilerplate/article/delivery/http/middleware"
	_articleRepo "github.com/khihadysucahyo/go-clean-arch-boilerplate/article/repository/mysql"
	_articleUcase "github.com/khihadysucahyo/go-clean-arch-boilerplate/article/usecase"
	_authorRepo "github.com/khihadysucahyo/go-clean-arch-boilerplate/author/repository/mysql"
)

func init() {
	viper.SetConfigFile(`.env`)
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if viper.GetBool(`DEBUG`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`DB_HOST`)
	dbPort := viper.GetString(`DB_PORT`)
	dbUser := viper.GetString(`DB_USER`)
	dbPass := viper.GetString(`DB_PASS`)
	dbName := viper.GetString(`DB_NAME`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	e.Use(middL.SENTRY)
	e.Use(middleware.Logger())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              viper.GetString(`SENTRY_DSN`),
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	e.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("APP_TIMEOUT")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	log.Fatal(e.Start(viper.GetString("APP_ADDRESS")))
}
