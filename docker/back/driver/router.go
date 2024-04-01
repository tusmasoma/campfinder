package driver

import (
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/tusmasoma/campfinder/docker/back/config"
	"github.com/tusmasoma/campfinder/docker/back/infra"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/handler"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/middleware"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

func InitRoute(serverConfig *config.ServerConfig) *chi.Mux {
	dialect := goqu.Dialect("mysql")
	db, err := config.NewDB()
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		return nil
	}
	client := config.NewClient()

	userRepo := infra.NewUserRepository(db)
	spotRepo := infra.NewSpotRepository(db)
	commentRepo := infra.NewCommentRepository(db, &dialect)
	imgRepo := infra.NewImageRepository(db, &dialect)
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
		MaxAge:           serverConfig.PreflightCacheDurationSec,
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

	return r
}
