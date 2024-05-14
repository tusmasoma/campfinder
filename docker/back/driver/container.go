package driver

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/dig"

	"github.com/tusmasoma/campfinder/docker/back/config"
	"github.com/tusmasoma/campfinder/docker/back/domain/repository"
	"github.com/tusmasoma/campfinder/docker/back/infra/mysql"
	"github.com/tusmasoma/campfinder/docker/back/infra/redis"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/handler"
	"github.com/tusmasoma/campfinder/docker/back/interfaces/middleware"
	"github.com/tusmasoma/campfinder/docker/back/usecase"
)

func BuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		providerSQLExecutor,
		config.NewClient,
		provideMySQLDialect,
		mysql.NewUserRepository,
		mysql.NewSpotRepository,
		mysql.NewCommentRepository,
		mysql.NewImageRepository,
		redis.NewSpotsRepository,
		redis.NewUserRepository,
		redis.NewCommentsRepository,
		redis.NewImagesRepository,
		usecase.NewUserUseCase,
		usecase.NewSpotUseCase,
		usecase.NewCommentUseCase,
		usecase.NewImageUseCase,
		usecase.NewAuthUseCase,
		handler.NewUserHandler,
		handler.NewSpotHandler,
		handler.NewCommentHandler,
		handler.NewImageHandler,
		middleware.NewAuthMiddleware,
		func(
			serverConfig *config.ServerConfig,
			userHandler handler.UserHandler,
			spotHandler handler.SpotHandler,
			commentHandler handler.CommentHandler,
			imgHandler handler.ImageHandler,
			authMiddleware middleware.AuthMiddleware,
		) *chi.Mux {
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
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}

	return container, nil
}

func provideMySQLDialect() *goqu.DialectWrapper {
	dialect := goqu.Dialect("mysql")
	return &dialect
}

func providerSQLExecutor() (repository.SQLExecutor, error) {
	return config.NewDB()
}
