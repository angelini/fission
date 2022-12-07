package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func StartServer(ctx context.Context, log *zap.Logger, svc string, hostAndPort string, handler http.Handler) {
	addr := hostAndPort
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf(":%s", hostAndPort)
	}

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	l := log.With(zap.String("service", svc), zap.String("addr", server.Addr))
	l.Info("starting server")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				l.Error("server error", zap.Error(err))
			}
		}
	}()
	<-ctx.Done()
	l.Info("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		if err != context.Canceled && err != context.DeadlineExceeded {
			l.Error("server shutdown error", zap.Error(err))
		}
	}
}
