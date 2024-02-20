package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tusmasoma/campfinder/docker/back/config"
	"github.com/tusmasoma/campfinder/docker/back/infra"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/handler"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/middleware"
	"github.com/tusmasoma/campfinder/docker/back/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	var addr string
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	Serve(addr)
}

func Serve(addr string) {
	db, err := config.NewDB()
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		return
	}
	client := config.NewClient()

	userRepo := infra.NewUserRepository(db)
	spotRepo := infra.NewSpotRepository(db)
	commentRepo := infra.NewCommentRepository(db)
	imgRepo := infra.NewImageRepository(db)
	redisRepo := infra.NewRedisRepository(client)

	userUseCase := usecase.NewUserUseCase(userRepo, redisRepo)
	spotUseCase := usecase.NewSpotUseCase(spotRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo)
	imgUseCase := usecase.NewImageUseCase(imgRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo)

	userHandler := handler.NewUserHandler(userUseCase, authUseCase)
	spotHandler := handler.NewSpotHandler(spotUseCase)
	commentHandler := handler.NewCommentHandler(commentUseCase, authUseCase)
	imgHandler := handler.NewImageHandler(imgUseCase, authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(redisRepo)

	/* ===== URLマッピングを行う ===== */
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link", "Authorization"},
		AllowCredentials: false,
		MaxAge:           config.PreflightCacheDurationSeconds,
	}))

	r.Use(middleware.Logging)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", userHandler.HandleUserCreate)
			r.Post("/login", userHandler.HandleUserLogin)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/api/user/logout", userHandler.HandleUserLogout)
			})
		})

		r.Route("/spot", func(r chi.Router) {
			r.Get("/", spotHandler.HandleSpotGet)
			r.Post("/create", spotHandler.HandleSpotCreate)
		})

		r.Route("/comment", func(r chi.Router) {
			r.Get("/", commentHandler.HandleCommentGet)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Post("/create", commentHandler.HandleCommentCreate)
				r.Post("/update", commentHandler.HandleCommentUpdate)
				r.Delete("/delete", commentHandler.HandleCommentDelete)
			})
		})

		r.Route("/img", func(r chi.Router) {
			r.Get("/", imgHandler.HandleImageGet)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Post("/create", imgHandler.HandleImageCreate)
				r.Post("/delete", imgHandler.HandleImageDelete)
			})
		})
	})
	/* ===== サーバの設定 ===== */
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}
	/* ===== サーバの起動 ===== */
	log.SetFlags(0)
	log.Println("Server running...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server stopping...")

	tctx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(tctx); err != nil {
		log.Println("failed to shutdown http server", err)
	}
	log.Println("Server exited")
}
