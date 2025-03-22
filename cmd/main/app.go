package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"server/internal/config"
	"server/internal/order"
	orderDB "server/internal/order/db"
	"server/internal/user"
	"server/internal/user/db"
	"server/pkg/client/mongodb"
	"server/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(),
		"mongodb+srv://maratmirzabalaev:15062004marat@cluster0.1egkm.mongodb.net",
		"", // Порт не нужен для SRV
		cfgMongo.Username,
		cfgMongo.Password,
		cfgMongo.Database,
		cfgMongo.AuthDB)

	if err != nil {
		panic(err)
	}

	fmt.Println("Config Collection =", cfg.MongoDB.Collection)
	userStorage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	userService := user.NewService(userStorage, logger)
	userHandler := user.NewHandler(logger, userService)

	fmt.Println("Config OrdersCollection =", cfg.MongoDB.OrdersCollection)

	orderStorage := orderDB.NewStorage(mongoDBClient, cfg.MongoDB.OrdersCollection, logger)
	orderService := order.NewService(orderStorage, logger)
	orderHandler := order.NewHandler(logger, orderService)

	logger.Info("register order handler ")
	orderHandler.Register(router)

	logger.Info("register user handler")
	userHandler.Register(router)

	// Middleware для игнорирования запросов к фронтенду
	ignoreFrontendRequests := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если запрос не начинается с /api, игнорируем его
			if !isAPIRequest(r.URL.Path) {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Применяем middleware
	handler := ignoreFrontendRequests(router)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handlerWithCORS := corsMiddleware.Handler(handler)
	start(&handlerWithCORS, cfg)
}

// Функция для проверки, является ли запрос API-запросом
func isAPIRequest(path string) bool {
	// Добавьте все ваши API-маршруты
	apiRoutes := []string{
		"/api/getOrders",
		"/api/createOrder",
		"/api/changeOrder",
		"/api/deleteOrder",
		"/api/users",
		"/api/user/:uuid",
		"/api/createuser",
		"/api/login",
	}

	for _, route := range apiRoutes {
		if path == route {
			return true
		}
	}
	return false
}

func start(handler *http.Handler, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("create unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening inix socket: %s", socketPath)

	} else {
		logger.Info("create tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port: %s: %s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      *handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
