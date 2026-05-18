// @title           MuchToDo API
// @version         1.0
// @description     This is an API for MuchToDo application with user authentication.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support - Innocent
// @contact.url    https://github.com/Innocent9712
// @contact.email  innocent@altschoolafrica.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and a JWT token."
package main

import (
        "context"
        "fmt"
        "log"
        "log/slog"
        "net/http"
        "os"
        "time"

        "github.com/gin-gonic/gin"
        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"

        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/auth"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/cache"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/config"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/database"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/handlers"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/logger"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/middleware"
        "github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/routes"

        _ "github.com/Innocent9712/much-to-do/Server/MuchToDo/docs" 
)

const usernameCacheSentinelKey = "username_cache_initialized"
const usernameCacheTTL = 24 * time.Hour

func main() {
        cfg, err := config.LoadConfig(".")
        if err != nil {
                log.Fatalf("could not load config: %v", err)
        }

        // BIND ENVIRONMENT OVERRIDES DIRECTLY TO BYPASS LOCAL FILE CHECKS
        envMongo := os.Getenv("MONGO_URI")
        if envMongo != "" {
                cfg.MongoURI = envMongo
        }

        logger.InitLogger(cfg)
        slog.Info("Logger initialized", "level", cfg.LogLevel, "format", cfg.LogFormat)

        dbClient, err := database.ConnectMongo(cfg.MongoURI, cfg.DBName)
        if err != nil {
                slog.Error("could not connect to MongoDB", slog.Any("error", err))
                os.Exit(1)
        }
        defer func() {
                if err = dbClient.Disconnect(context.Background()); err != nil {
                        slog.Error("Error disconnecting from MongoDB", slog.Any("error", err))
                }
        }()
        slog.Info("Successfully connected to MongoDB.")

        cacheService := cache.NewCacheService(cfg)
        tokenService := auth.NewTokenService(cfg.JWTSecretKey, cfg.JWTExpirationHours)

        preloadUsernamesIntoCache(dbClient, cacheService, cfg)

        router := setupRouter(dbClient, cfg, tokenService, cacheService)

        // FORCE THE RUNNER ENGINE TO BIND TO ALL INTERFACES ON PORT 8080
        address := "0.0.0.0:8080"
        slog.Info("🚀 StartTech Server binding globally", "address", address)
        
        if err := http.ListenAndServe(address, router); err != nil {
                log.Fatalf("❌ Critical server boot crash: %v", err)
        }
}

func preloadUsernamesIntoCache(db *mongo.Client, cacheSvc cache.Cache, cfg config.Config) {
        if !cfg.EnableCache {
                slog.Info("Caching is disabled. Skipping username preloading.")
                return
        }
        var sentinelVal string
        err := cacheSvc.Get(context.Background(), usernameCacheSentinelKey, &sentinelVal)
        if err == nil {
                slog.Info("Username cache already initialized. Skipping preload.")
                return
        }
        slog.Info("Preloading usernames into cache...")
        ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
        defer cancel()
        userCollection := db.Database(cfg.DBName).Collection("users")
        opts := options.Find().SetProjection(bson.M{"username": 1})
        cursor, err := userCollection.Find(ctx, bson.D{}, opts)
        if err != nil {
                slog.Error("Error querying for usernames to preload", slog.Any("error", err))
                return
        }
        defer cursor.Close(ctx)
        usernamesToCache := make(map[string]interface{})
        for cursor.Next(ctx) {
                var result struct {
                        Username string `bson:"username"`
                }
                if err := cursor.Decode(&result); err != nil {
                        slog.Warn("Error decoding username during preload", slog.Any("error", err))
                        continue
                }
                if result.Username != "" {
                        cacheKey := fmt.Sprintf("username-taken:%s", result.Username)
                        usernamesToCache[cacheKey] = true
                }
        }
        if len(usernamesToCache) > 0 {
                cacheSvc.SetMany(ctx, usernamesToCache, usernameCacheTTL)
                cacheSvc.Set(ctx, usernameCacheSentinelKey, "true", usernameCacheTTL)
        }
}

func setupRouter(db *mongo.Client, cfg config.Config, tokenSvc *auth.TokenService, cacheSvc cache.Cache) *gin.Engine {
        gin.SetMode(gin.ReleaseMode)
        router := gin.Default()
        todoCollection := db.Database(cfg.DBName).Collection("todos")
        userCollection := db.Database(cfg.DBName).Collection("users")
        todoHandler := handlers.NewTodoHandler(todoCollection)
        userHandler := handlers.NewUserHandler(userCollection, todoCollection, tokenSvc, cacheSvc, db, cfg)
        healthHandler := handlers.NewHealthHandler(db, cacheSvc, cfg.EnableCache)
        corsMiddleware := middleware.CORSMiddleware(cfg.AllowedOrigins)
        authMiddleware := middleware.AuthMiddleware(tokenSvc, cfg)
        router.Use(corsMiddleware)
        routes.RegisterRoutes(router, userHandler, todoHandler, healthHandler, authMiddleware)
        router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
        router.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Welcome to MuchToDo API"}) })
        router.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"}) })
        return router
}
