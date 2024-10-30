package transport

type authContextKey struct{}

// func AuthUnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, status.Error(codes.Unauthenticated, "missing metadata")
// 	}

// 	if authHeader, ok := md["autorization"]; ok && len(authHeader) > 0 {
// 		token := authHeader[0]

// 		user, err := verifyToken(token)

// 		if err != nil {
// 			return nil, status.Error(codes.Unauthenticated, "invalid token")
// 		}

// 		ctx = context.WithValue(ctx, authContextKey{}, user)
// 	}

// 	return handler(ctx, req)
// }
