package main

import (
	"context"
	"fmt"
)

type userIdKey struct{}
type authTokenKey struct{}

func main() {
	processRequest("takeshi", "hoobar")
}

func processRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), userIdKey{}, userId)
	ctx = context.WithValue(ctx, authTokenKey{}, authToken)
	handleResponse(ctx)
}

func handleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (%v)\n",
		ctx.Value(userIdKey{}),
		ctx.Value(authTokenKey{}),
	)
}
