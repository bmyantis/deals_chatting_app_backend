package main

import (
	"context"
	"fmt"
    "strings"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
    "go.opentelemetry.io/otel/semconv/v1.4.0"
    
	"github.com/go-playground/validator/v10"
	gocloak "github.com/Nerzal/gocloak/v13"
    
	"deals_chatting_app_backend/internal/controller"
	// "deals_chatting_app_backend/internal/middleware"
	"deals_chatting_app_backend/internal/model"
	"deals_chatting_app_backend/internal/repository"
	"deals_chatting_app_backend/internal/service"
	"deals_chatting_app_backend/internal/config"
	"deals_chatting_app_backend/internal/database"
	"deals_chatting_app_backend/internal/router"

)

func main() {
	fmt.Println("Hello from myapp")

	// Initialize Logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Initialize Config
	config.InitConfig(logger)

	// // Initialize Tracer
	// shutdown, err := initTracer()
	// if err != nil {
	// 	logger.Sugar().Panicf("failed to initialize trace provider: %v", err)
	// }
	// defer func() {
	// 	if err := shutdown(context.Background()); err != nil {
	// 		logger.Sugar().Warnf("Error shutting down tracer provider: %v", err)
	// 	}
	// }()

	db := database.DatabaseConnection()
	if viper.GetBool("autoMigrate") {
		zap.L().Sugar().Infof("Executing autoMigrate")
		db.AutoMigrate(
			&model.User{},
			&model.Profile{},
			&model.Preferences{},
            &model.Swipe{},
		)
	}
    

	// Keycloak
    // keycloakURL := strings.TrimSuffix(viper.GetString("KEYCLOAK_URL"), "/")
    // fmt.Printf("Keycloak URL: %s\n", keycloakURL) // Add this log
	// keycloak := gocloak.NewClient(keycloakURL)
    // keycloak.RestyClient().SetDebug((strings.ToUpper(viper.GetString("KEYCLOAK_DEBUG")) == "TRUE"))
	// Keycloak
	keycloak := gocloak.NewClient(viper.GetString("KEYCLOAK_URL"))
	keycloak.RestyClient().SetDebug((strings.ToUpper(viper.GetString("KEYCLOAK_DEBUG")) == "TRUE"))

    
	validator := validator.New()

	// Repositories
	userRepository := repository.NewUserRepository(db)
    swipeRepository := repository.NewSwipeRepository(db)

	// Services
	userService := service.NewUserService(userRepository, keycloak)
    swipeService := service.NewSwipeService(swipeRepository)

	// Controllers
    userController := controller.NewUserController(userService, validator)
	swipeController := controller.NewSwipeController(swipeService, validator)	

	// Create a new Gin router instance by calling NewRouter function
	r := router.NewRouter(keycloak, userController, swipeController, logger) // Use the router instance returned by NewRouter

	// Middlewares
	// r.Use(middleware.LoggerMiddleware())
	// r.Use(middleware.TraceAndSpanMiddleware())

	// Log all routes for debugging purposes
	for _, route := range r.Routes() {
		fmt.Printf("Route: %s %s\n", route.Method, route.Path)
	}

	// Start Server
	port := viper.GetString("http_port")
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Sugar().Fatalf("failed to run server: %v", err)
	}
}

func initTracer() (func(ctx context.Context) error, error) {
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("my-service"),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
