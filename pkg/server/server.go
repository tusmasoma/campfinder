package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tusmasoma/campfinder/cache"
	"github.com/tusmasoma/campfinder/db"
	"github.com/tusmasoma/campfinder/internal/auth"
	"github.com/tusmasoma/campfinder/pkg/http/middleware"
	"github.com/tusmasoma/campfinder/pkg/server/handler"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

const (
	readTimeout             = 5 * time.Second
	writeTimeout            = 10 * time.Second
	idleTimeout             = 15 * time.Second
	gracefulShutdownTimeout = 5 * time.Second
)

var (
	dbUser        = os.Getenv("MYSQL_ROOT_USER")
	dbPassword    = os.Getenv("MYSQL_ROOT_PASSWORD")
	dbHost        = os.Getenv("MYSQL_HOST")
	dbPort        = os.Getenv("MYSQL_PORT")
	dbName        = os.Getenv("MYSQL_DB_NAME")
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDB, _    = strconv.Atoi(os.Getenv("REDIS_DB"))
)

func Serve(addr string) {
	var err error

	c := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Println(c)
	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		return
	}
	client := redis.NewClient(&redis.Options{Addr: redisAddr, Password: redisPassword, DB: redisDB})

	userRepo := db.NewUserRepository(database)
	spotRepo := db.NewSpotRepository(database)
	commentRepo := db.NewCommentRepository(database)
	imgRepo := db.NewImageRepository(database)
	redisRepo := cache.NewRedisRepository(client)
	authMiddleware := middleware.NewAuthMiddleware(redisRepo)
	authHandler := auth.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo, redisRepo, authHandler)
	spotHandler := handler.NewSpotHandler(spotRepo)
	commentHandler := handler.NewCommentHandler(commentRepo, authHandler)
	imgHandler := handler.NewImageHandler(imgRepo, authHandler)

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
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
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

	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
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
