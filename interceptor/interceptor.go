package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func toSet(arr []string) map[string]interface{} {
	m := make(map[string]interface{})
	for _, v := range arr {
		m[v] = nil
	}
	return m
}

// StreamBlacklistConditionalInterceptor will only invoke the interceptor if the method ISN'T in the blacklist.
func StreamBlacklistConditionalInterceptor(methods []string, interceptor grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	set := toSet(methods)
	fn := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := set[info.FullMethod]; ok == true {
			return handler(srv, ss)
		}
		return interceptor(srv, ss, info, handler)
	}
	return fn
}

// StreamWhitelistConditionalInterceptor will only invoke the interceptor if the method IS in the whitelist.
func StreamWhitelistConditionalInterceptor(methods []string, interceptor grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	set := toSet(methods)
	fn := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := set[info.FullMethod]; ok == false {
			return handler(srv, ss)
		}
		return interceptor(srv, ss, info, handler)
	}
	return fn
}

// UnaryBlacklistConditionalInterceptor will only invoke the interceptor if the method ISN'T in the blacklist.
func UnaryBlacklistConditionalInterceptor(methods []string, interceptor grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	set := toSet(methods)
	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := set[info.FullMethod]; ok == true {
			return handler(ctx, req)
		}
		return interceptor(ctx, req, info, handler)
	}
	return fn
}

// UnaryWhitelistConditionalInterceptor will only invoke the interceptor if the method IS in the whitelist.
func UnaryWhitelistConditionalInterceptor(methods []string, interceptor grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	set := toSet(methods)
	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := set[info.FullMethod]; ok == false {
			return handler(ctx, req)
		}
		return interceptor(ctx, req, info, handler)
	}
	return fn
}
