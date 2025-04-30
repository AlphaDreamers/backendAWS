package server

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/auth/cmd/config"
	"github.com/SwanHtetAungPhyo/auth/internal/delivery"
	"github.com/SwanHtetAungPhyo/auth/internal/models"
	"github.com/SwanHtetAungPhyo/auth/internal/repo"
	"github.com/SwanHtetAungPhyo/auth/internal/service"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"runtime"
	"time"
)

type App interface {
	Start()
	Stop()
	setupMiddlewares()
	setupRoutes()
	setupDatabase() error
	setupRabbitMq() error
	setupRedis() error
	setupLogger()
	migration()
	cleanupDatabase()
	cleanupRabbitMq()
	cleanupRedis()
}

type AppState struct {
	app            *fiber.App
	database       *gorm.DB
	log            *logrus.Logger
	redisClient    *redis.Client
	config         *config.Config
	rabbitMqClient *amqp091.Connection
	delivery       delivery.Delivery
}

func NewApp() App {
	app := AppState{}
	app.setupLogger()
	userRepo := repo.NewUserRepository(app.database)
	app.delivery = delivery.NewUserDelivery(app.log, service.NewUserService(userRepo))

	if err := app.setupRedis(); err != nil {
		app.log.Warnf("Failed to setup Redis: ", err)
	}

	if err := app.setupRabbitMq(); err != nil {
		app.log.Warnf("Failed to setup RabbitMQ: ", err)
	}

	if err := app.setupDatabase(); err != nil {
		app.log.Warnf("Failed to setup Database: ", err)
	}

	return &app
}

func (a *AppState) Start() {
	a.app = fiber.New(fiber.Config{
		ReadTimeout:           10 * 60,
		WriteTimeout:          10 * 60,
		WriteBufferSize:       1024 * 1024,
		ReadBufferSize:        1024 * 1024,
		IdleTimeout:           10 * 60,
		Concurrency:           1000,
		DisableStartupMessage: false,
	})
	a.migration()
	a.setupMiddlewares()
	a.setupRoutes()

	a.log.Info("application started....")
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	err := a.app.Listen(":" + port)
	if err != nil {
		a.log.Fatal("Failed to start server: ", err)
	}
}

func (a *AppState) Stop() {
	a.log.Info("application stopped....")
	a.cleanupDatabase()
	a.cleanupRabbitMq()
	a.cleanupRedis()
	if err := a.app.Shutdown(); err != nil {
		a.log.Error("Error while shutting down: ", err)
	}
}

func (a *AppState) setupMiddlewares() {
	// Add middleware setup if needed
}

func (a *AppState) setupRoutes() {
	a.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	authRoutes := a.app.Group("/auth")
	authRoutes.Post("/login", a.delivery.Login)
	authRoutes.Post("/register", a.delivery.Register)
	a.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	authRoutes.Post("/logout", a.delivery.Logout)
	authRoutes.Post("/forgot-password", a.delivery.ForgotPassword)
	authRoutes.Post("/refresh", a.delivery.Refresh)
	authRoutes.Post("/reset-password", a.delivery.ResetPassword)
}

func (a *AppState) setupDatabase() error {
	dsn := viper.GetString("database.dsn")
	for i := 0; i < 5; i++ {
		var err error
		a.database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		a.log.Errorf("Database connection failed (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if a.database == nil {
		return fmt.Errorf("failed to connect to database after 5 attempts")
	}

	sqlDB, err := a.database.DB()
	if err != nil {
		return fmt.Errorf("failed to set up database connection pool: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(25 * time.Minute)
	a.log.Info("Database connected successfully!")
	return nil
}

func (a *AppState) setupRabbitMq() error {
	for i := 0; i < 5; i++ {
		var err error
		a.rabbitMqClient, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
		if err == nil {
			break
		}
		a.log.Errorf("RabbitMQ connection failed (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if a.rabbitMqClient == nil {
		return fmt.Errorf("failed to connect to RabbitMQ after 5 attempts")
	}

	a.log.Info("RabbitMQ connected successfully!")
	return nil
}

func (a *AppState) GetRabbitMqClient() *amqp091.Connection {
	return a.rabbitMqClient
}

func (a *AppState) GetRedisClient() *redis.Client {
	return a.redisClient
}

func (a *AppState) migration() {
	err := a.database.Migrator().AutoMigrate(models.UserInDB{})
	if err != nil {
		a.log.Error("Failed to migrate database: ", err.Error())
		return
	}
}

func (a *AppState) setupRedis() error {
	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err := a.redisClient.Ping(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}
	a.log.Info("Redis connected successfully!")
	return nil
}

func (a *AppState) setupLogger() {
	a.log = logrus.New()
	a.log.SetLevel(logrus.DebugLevel)
	a.log.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			if f.Function == "" || f.File == "" {
				return "unknown_file", "unknown_function"
			}
			return f.File, f.Function
		},
	})

}

func (a *AppState) cleanupDatabase() {
	if a.database != nil {
		db, err := a.database.DB()
		if err != nil {
			a.log.Error("Failed to get DB instance: ", err)
			return
		}
		err = db.Close()
		if err != nil {
			a.log.Error("Error closing DB connection: ", err)
		}
	}
}

func (a *AppState) cleanupRabbitMq() {
	if a.rabbitMqClient != nil {
		err := a.rabbitMqClient.Close()
		if err != nil {
			a.log.Error("Error closing RabbitMQ connection: ", err)
		}
	}
}

func (a *AppState) cleanupRedis() {
	if a.redisClient != nil {
		err := a.redisClient.Close()
		if err != nil {
			a.log.Error("Error closing Redis connection: ", err)
		}
	}
}
