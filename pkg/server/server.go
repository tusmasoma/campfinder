package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tusmasoma/campfinder/cache"
	"github.com/tusmasoma/campfinder/db"
	"github.com/tusmasoma/campfinder/internal/auth"
	"github.com/tusmasoma/campfinder/pkg/http/middleware"
	"github.com/tusmasoma/campfinder/pkg/server/handler"

	_ "github.com/go-sql-driver/mysql"
)

func Serve(addr string) {
	var err error

	DB, err := sql.Open("mysql", "root:campfinder@tcp(mysql:3306)/campfinderdb")
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "", DB: 0})

	userRepo := db.NewUserRepository(DB)
	spotRepo := db.NewSpotRepository(DB)
	commentRepo := db.NewCommentRepository(DB)
	redisRepo := cache.NewRedisRepository(client)
	authMiddleware := middleware.NewAuthMiddleware(redisRepo)
	authHandler := auth.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo, redisRepo, authHandler)
	spotHandler := handler.NewSpotHandler(spotRepo)
	commentHandler := handler.NewCommentHandler(commentRepo, authHandler)

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/api/user/create", middleware.Logging(post(userHandler.HandleUserCreate)))
	http.HandleFunc("/api/user/login", middleware.Logging(post(userHandler.HandleUserLogin)))
	http.HandleFunc("/api/user/logout", middleware.Logging(get(authMiddleware.Authenticate(userHandler.HandleUserLogout))))
	http.HandleFunc("/api/spot", middleware.Logging(get(spotHandler.HandleSpotGet)))
	http.HandleFunc("/api/spot/create", middleware.Logging(post(spotHandler.HandleSpotCreate)))
	http.HandleFunc("/api/comment", middleware.Logging(get(commentHandler.HandleCommentGet)))
	http.HandleFunc("/api/comment/create", middleware.Logging(post(authMiddleware.Authenticate(commentHandler.HandleCommentCreate))))
	http.HandleFunc("/api/comment/update", middleware.Logging(post(authMiddleware.Authenticate(commentHandler.HandleCommentUpdate))))
	http.HandleFunc("/api/comment/delete", middleware.Logging(post(authMiddleware.Authenticate(commentHandler.HandleCommentDelete))))

	/* ===== サーバの設定 ===== */
	srv := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	/* ===== サーバの起動 ===== */
	log.SetFlags(0)
	log.Println("Server running...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
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
