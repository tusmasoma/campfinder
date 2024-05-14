package driver

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tusmasoma/campfinder/docker/back/config"
)

func Run() {
	var addr string
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	serverConfig, err := config.NewServerConfig(context.Background())
	if err != nil {
		log.Printf("Failed to load server config: %s\n", err)
		return
	}
	/* ===== サーバの設定 ===== */
	srv := &http.Server{
		Addr:         addr,
		Handler:      InitRoute(serverConfig),
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
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

	tctx, cancel := context.WithTimeout(context.Background(), serverConfig.GracefulShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(tctx); err != nil {
		log.Println("failed to shutdown http server", err)
	}
	log.Println("Server exited")
}
