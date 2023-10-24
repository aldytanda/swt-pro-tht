package main

import (
	"os"

	"github.com/aldytanda/swt-pro-tht/generated"
	"github.com/aldytanda/swt-pro-tht/handler"
	"github.com/aldytanda/swt-pro-tht/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn:          dbDsn,
		JWTSecretKey: jwtSecret,
	})
	opts := handler.NewServerOptions{
		JWTSecretKey: jwtSecret,
		Repository:   repo,
	}
	return handler.NewServer(opts)
}
