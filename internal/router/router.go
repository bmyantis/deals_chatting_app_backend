package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	gocloak "github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"

	"deals_chatting_app_backend/internal/controller"
	"deals_chatting_app_backend/internal/middleware"
)

func NewRouter(keycloak *gocloak.GoCloak, userController controller.UserController, swipeController controller.SwipeController, logger *zap.Logger) *gin.Engine {
	keycloakClientId := viper.GetString("KEYCLOAK_CLIENT_ID")
	keycloakRealm := viper.GetString("KEYCLOAK_REALM")
	keycloakClientSecret := viper.GetString("KEYCLOAK_CLIENT_SECRET")

	router := gin.Default()
	router.Use(middleware.Middleware(logger))
	router.Use(middleware.PaginationMiddleware())

	// Add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK!")
	})
	router.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	})

	v1Router := router.Group("/v1")
	userRouter := v1Router.Group("/user")
	userRouter.POST("/signup", userController.Signup)
	userRouter.POST("/login", userController.Login)

	// Apply Keycloak middleware to routes that require authentication
	authenticatedUser := userRouter.Group("/")
	authenticatedUser.Use(middleware.KeycloakAuthMiddleware(keycloak, keycloakClientId, keycloakClientSecret, keycloakRealm))
	authenticatedUser.PUT("/:id/profile", userController.CreateOrUpdateProfile)
	authenticatedUser.PUT("/:id/preferences", userController.CreateOrUpdatePreferences)
	authenticatedUser.GET("/", userController.FindAll)

	swipeRouter := v1Router.Group("/swipe")
	authenticatedSwipe := swipeRouter.Group("/")
	authenticatedSwipe.Use(middleware.KeycloakAuthMiddleware(keycloak, keycloakClientId, keycloakClientSecret, keycloakRealm))
	authenticatedSwipe.POST("/", swipeController.CreateSwipe)

	return router
}
