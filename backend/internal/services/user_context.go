package services

import "context"

type CurrentUserCtxKeyType struct {
}

var CurrentUserCtxKey = CurrentUserCtxKeyType{}

func PutUsernameInContex(iCtx context.Context, iUsername string) context.Context {
	return context.WithValue(iCtx, CurrentUserCtxKey, iUsername)
}

func GetCurrentUserFromContext(ctx context.Context) (int, error) {
	/// for temporary testing
	return 1, nil
	// if raw, ok := ctx.Value(UserCtxKey).(string); ok {
	// 	return raw
	// } else {
	// 	return ""
	// }
}
