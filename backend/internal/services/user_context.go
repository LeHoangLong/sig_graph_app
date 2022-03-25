package services

import "context"

var UserCtxKey = "username"

func PutUsernameInContex(iCtx context.Context, iUsername string) context.Context {
	return context.WithValue(iCtx, UserCtxKey, iUsername)
}

func GetUserFromContext(ctx context.Context) int {
	/// for temporary testing
	return 1
	// if raw, ok := ctx.Value(UserCtxKey).(string); ok {
	// 	return raw
	// } else {
	// 	return ""
	// }
}
