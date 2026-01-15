package main

import (
	"context"
	"os"
	"os/signal"
	"restapi-app"
	grpc_auth "restapi-app/pkg/client/auth/grpc"
	"restapi-app/pkg/handler"
	"restapi-app/pkg/repository"
	"restapi-app/pkg/service"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const level string = "local"

func main() {
	//logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := InitConfig(); err != nil {
		logrus.Fatalf("Cant read config file: %s", err)
	}
	_ = godotenv.Load() // сделал так, чтобы докер не ругался
	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading env vars: %s", err)
	// }

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("cant connect to db %s", err)
	}

	// init grpc Client
	authClient, err := grpc_auth.NewClient(context.Background(), viper.GetString("client.address"))
	if err != nil {
		logrus.Fatalf("Error conecting to auth grpc server: %s", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, authClient)

	srv := new(restapi.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error running server: %v", err)
		}
	}()

	logrus.Print("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App is shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Error occured on server shutting down: %s", err.Error())
	}
	logrus.Print("server shutdown")

	if err = db.Close(); err != nil {
		logrus.Errorf("Error occured on db connection close: %s", err.Error())
	}
	logrus.Print("db closed")
}

func InitConfig() error {
	switch level {
	case "local":
		viper.AddConfigPath("configs")
	case "docker":
		viper.AddConfigPath(".")
	}

	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
