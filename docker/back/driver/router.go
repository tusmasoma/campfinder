package driver

import (
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/tusmasoma/campfinder/docker/back/config"
	"github.com/tusmasoma/campfinder/docker/back/infra/mysql"
	"github.com/tusmasoma/campfinder/docker/back/infra/redis"
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

	userRepo := mysql.NewUserRepository(db, &dialect)
	spotRepo := mysql.NewSpotRepository(db, &dialect)
	commentRepo := mysql.NewCommentRepository(db, &dialect)
	imgRepo := mysql.NewImageRepository(db, &dialect)

	spotsRedisRepo := redis.NewSpotsRepository(client)
	userRedisRepo := redis.NewUserRepository(client)
	imgRedisRepo := redis.NewImagesRepository(client)

	userUseCase := usecase.NewUserUseCase(userRepo, userRedisRepo)
	spotUseCase := usecase.NewSpotUseCase(spotRepo, spotsRedisRepo)
	commentUseCase := usecase.NewCommentUseCase(commentRepo)
	imgUseCase := usecase.NewImageUseCase(imgRepo, imgRedisRepo)
	authUseCase := usecase.NewAuthUseCase(userRepo)

	userHandler := handler.NewUserHandler(userUseCase, authUseCase)
	spotHandler := handler.NewSpotHandler(spotUseCase)
	commentHandler := handler.NewCommentHandler(commentUseCase, authUseCase)
	imgHandler := handler.NewImageHandler(imgUseCase, authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(userRedisRepo)

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
			r.Post("/create", userHandler.CreateUser)
			r.Post("/login", userHandler.Login)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/api/user/logout", userHandler.Logout)
			})
		})

		r.Route("/spot", func(r chi.Router) {
			r.Get("/", spotHandler.ListSpots)
			r.Post("/create", spotHandler.CreateSpot)
		})

		r.Route("/comment", func(r chi.Router) {
			r.Get("/", commentHandler.ListComments)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Post("/create", commentHandler.CreateComment)
				r.Post("/update", commentHandler.UpdateComment)
				r.Delete("/delete", commentHandler.DeleteComment)
			})
		})

		r.Route("/img", func(r chi.Router) {
			r.Get("/", imgHandler.ListImages)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Post("/create", imgHandler.CreateImage)
				r.Post("/delete", imgHandler.DeleteImage)
			})
		})
	})

	return r
}
