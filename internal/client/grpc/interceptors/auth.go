// Package interceptors содержит gRPC интерцепторы, которые используются для добавления аутентификационных данных в метаданные запросов.
package interceptors

import (
	"context"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// AddAuth возвращает UnaryClientInterceptor, который добавляет токен доступа и идентификатор клиента в метаданные запроса.
// Если токен пуст, вызов переходит к следующему обработчику без изменения контекста.
func AddAuth(token *string, clientID uint32) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if len(*token) == 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		md := metadata.New(map[string]string{
			consts.AccessTokenHeader: *token,
			consts.ClientIDHeader:    strconv.Itoa(int(clientID)),
		})

		mdCtx := metadata.NewOutgoingContext(ctx, md)
		return invoker(mdCtx, method, req, reply, cc, opts...)
	}
}
