package update

import (
	"authentication-app/storage/db"
	"context"
)

type TokenUpdater interface {
	UpdateToken(ctx context.Context, arg *db.UpdateTokenParams) (*db.User, error)
}
