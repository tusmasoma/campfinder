package driver

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/dig"

	"github.com/tusmasoma/campfinder/docker/back/config"
	"github.com/tusmasoma/campfinder/docker/back/infra/mysql"
	"github.com/tusmasoma/campfinder/docker/back/infra/redis"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/handler"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/middleware"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(config.NewServerConfig)
	container.Provide(config.NewDB)
	container.Provide(config.NewClient)
	container.Provide(goqu.Dialect("mysql"))

	container.Provide(mysql.NewUserRepository)
	container.Provide(mysql.NewSpotRepository)
	container.Provide(mysql.NewCommentRepository)
	container.Provide(mysql.NewImageRepository)
	container.Provide(redis.NewSpotsRepository)
	container.Provide(redis.NewUserRepository)
	container.Provide(redis.NewCommentsRepository)
	container.Provide(redis.NewImagesRepository)

	container.Provide(usecase.NewUserUseCase)
	container.Provide(usecase.NewSpotUseCase)
	container.Provide(usecase.NewCommentUseCase)
	container.Provide(usecase.NewImageUseCase)
	container.Provide(usecase.NewAuthUseCase)

	container.Provide(handler.NewUserHandler)
	container.Provide(handler.NewSpotHandler)
	container.Provide(handler.NewCommentHandler)
	container.Provide(handler.NewImageHandler)

	container.Provide(middleware.NewAuthMiddleware)

	container.Provide(func(serverConfig config.ServerConfig, userHandler handler.UserHandler, spotHandler handler.SpotHandler, commentHandler handler.CommentHandler, imgHandler handler.ImageHandler, authMiddleware middleware.AuthMiddleware) *chi.Mux {
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
				r.Get("/{spotID}", spotHandler.GetSpot)
				r.Post("/create", spotHandler.CreateSpot)
				r.Post("/batchcreate", spotHandler.BatchCreateSpots)
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
	})

	return container
}
