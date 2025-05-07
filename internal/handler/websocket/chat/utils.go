package chat

import (
	"context"

	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/uptrace/bun"
)

func checkAuthorization(db *bun.DB, userID int64) bool {
	user := &user.User{ID: userID}
	err := db.NewSelect().Model(user).Where("id = ?", userID).Scan(context.Background())
	return err == nil
}
