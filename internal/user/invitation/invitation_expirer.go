package invitation

import (
	"context"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/pkg/sqldb"
)

type expirer struct {
	getterRepo  GetterRepo
	updaterRepo UpdaterRepo
}

func (e *expirer) ExpireInvitationsByEmailTx(ctx context.Context, tx sqldb.Executable, email string) error {
	//TODO implement me
	panic("implement me")
}

type GetterRepo interface {
	ListByEmailTx(ctx context.Context, tx sqldb.Queryable, email string) ([]entity.UserInvitation, error)
}

type UpdaterRepo interface {
	UpdateInvitationTx(
		ctx context.Context, tx sqldb.Executable, invitation entity.UserInvitation,
	) (entity.UserInvitation, error)
}
