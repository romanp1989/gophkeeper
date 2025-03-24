// Package interceptors предоставляет набор gRPC интерцепторов для клиентской стороны, которые могут быть использованы для управления поведением вызовов, таких как установка таймаутов на операции.
package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// Timeout возвращает UnaryClientInterceptor, который устанавливает таймаут для gRPC вызовов.
// Интерцептор модифицирует контекст вызова, добавляя к нему ограничение по времени.
func Timeout(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timedCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(timedCtx, method, req, reply, cc, opts...)
	}
}
