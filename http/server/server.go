package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/ivanov-nikolay/game/http/server/handler"
	"github.com/ivanov-nikolay/game/internal/service"
)

// Маршрутизация
func newMux(ctx context.Context, logger *zap.Logger, height, width int, lifeService service.LifeService) (http.Handler, error) {
	muxHandler, err := handler.New(ctx, height, width, lifeService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}
	// Middleware для обработчиков
	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

// Run ...
func Run(ctx context.Context, logger *zap.Logger, height, width, fill int) (func(context.Context) error, error) {
	// Сервис с игрой
	lifeService, err := service.New(height, width, fill)
	if err != nil {
		return nil, err
	}

	muxHandler, err := newMux(ctx, logger, height, width, *lifeService) //todo
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: ":8081", Handler: muxHandler}

	go func() {
		// Запускаем сервер
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()
	// Вернём функцию для завершения работы сервера
	return srv.Shutdown, nil
}

// Middleware для логированя запросов
func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропуск запроса к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершение логирования после выполнения запроса
			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
