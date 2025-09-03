package utils

import (
	"SaltySpitoon/internal/constants"
	"context"
)

func GetUserIDFromCtx(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(constants.UserIDCtxKey).(int64)
	return userID, ok
}
