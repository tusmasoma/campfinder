package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tusmasoma/campfinder/config"
	"github.com/tusmasoma/campfinder/infra"
	"github.com/tusmasoma/campfinder/interfaces/handler"
	"github.com/tusmasoma/campfinder/interfaces/middleware"
	"github.com/tusmasoma/campfinder/usecase"
)

func Serve(addr string) {
	var err error

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
	http.HandleFunc("/api/user/create",
		middleware.Logging(post(userHandler.HandleUserCreate)))
	http.HandleFunc("/api/user/login",
		middleware.Logging(post(userHandler.HandleUserLogin)))
	http.HandleFunc("/api/user/logout",
		middleware.Logging(get(authMiddleware.Authenticate(userHandler.HandleUserLogout))))
	http.HandleFunc("/api/spot",
		middleware.Logging(get(spotHandler.HandleSpotGet)))
	http.HandleFunc("/api/spot/create",
		middleware.Logging(post(spotHandler.HandleSpotCreate)))
	http.HandleFunc("/api/comment",
		middleware.Logging(get(commentHandler.HandleCommentGet)))
	http.HandleFunc("/api/comment/create",
		middleware.Logging(post(authMiddleware.Authenticate(commentHandler.HandleCommentCreate))))
	http.HandleFunc("/api/comment/update",
		middleware.Logging(post(authMiddleware.Authenticate(commentHandler.HandleCommentUpdate))))
	http.HandleFunc("/api/comment/delete",
		middleware.Logging(del(authMiddleware.Authenticate(commentHandler.HandleCommentDelete))))
	http.HandleFunc("/api/img",
		middleware.Logging(get(imgHandler.HandleImageGet)))
	http.HandleFunc("/api/img/create",
		middleware.Logging(post(authMiddleware.Authenticate(imgHandler.HandleImageCreate))))
	http.HandleFunc("/api/img/delete",
		middleware.Logging(del(authMiddleware.Authenticate(imgHandler.HandleImageDelete))))

	/* ===== サーバの設定 ===== */
	srv := &http.Server{
		Addr:         addr,
		Handler:      nil,
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

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// delete DELETEリクエストを処理する
func del(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodDelete)
}

func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS対応
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token,Authorization")
		w.Header().Add("Access-Control-Expose-Headers", "Authorization")

		// プリフライトリクエストは処理を通さない
		if r.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			if _, err := w.Write([]byte("Method Not Allowed")); err != nil {
				log.Printf("Error writing data: %v", err)
			}
			return
		}

		// 共通のレスポンスヘッダを設定
		w.Header().Add("Content-Type", "application/json")
		apiFunc(w, r)
	}
}
