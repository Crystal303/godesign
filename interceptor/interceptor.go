package main

import (
	"context"
	"fmt"
)

type interceptor func(c context.Context, ivk invoker) error

type invoker func(c context.Context, interceptors []interceptor) error

func h1(c context.Context) {
	fmt.Println("h1 do something")
}
func h2(c context.Context) {
	fmt.Println("h2 do something")
}
func h3(c context.Context) {
	fmt.Println("h3 do something")
}

func main() {
	ceps := make([]interceptor, 0, 4)
	ceps = append(ceps, func(c context.Context, ivk invoker) error {
		h1(c)
		return ivk(c, ceps)
	})
	ceps = append(ceps, func(c context.Context, ivk invoker) error {
		h2(c)
		return ivk(c, ceps)
	})
	ceps = append(ceps, func(c context.Context, ivk invoker) error {
		h3(c)
		return ivk(c, ceps)
	})
	var ivk invoker = func(c context.Context, interceptors []interceptor) error {
		fmt.Println("invoker")
		return nil
	}
	cep := chainUnaryInterceptors(ceps)

	ctx := context.Background()
	_ = cep(ctx, ivk)
}

func chainUnaryInterceptors(interceptors []interceptor) interceptor {
	if len(interceptors) == 0 {
		return nil
	}
	if len(interceptors) == 1 {
		return interceptors[0]
	}
	return func(ctx context.Context, ivk invoker) error {
		return interceptors[0](ctx, getChainUnaryInvoker(interceptors, 0, ivk))
	}
}

func getChainUnaryInvoker(interceptors []interceptor, curr int, finalInvoker invoker) invoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(c context.Context, interceptors []interceptor) error {
		return interceptors[curr+1](c, getChainUnaryInvoker(interceptors, curr+1, finalInvoker))
	}
}
