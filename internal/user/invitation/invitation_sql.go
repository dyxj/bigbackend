package invitation

import (
	"time"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/google/uuid"
)

const (
	dbcUkToken                = "user_invitation_token_uk"
	dbcUkAcceptedPendingEmail = "user_invitation_accepted_pending_email_uk"
)

type userInvitationAuditableEntity struct{ E *entity.UserInvitation }

func (a userInvitationAuditableEntity) GetID() uuid.UUID          { return a.E.ID }
func (a userInvitationAuditableEntity) SetID(id uuid.UUID)        { a.E.ID = id }
func (a userInvitationAuditableEntity) SetCreateTime(t time.Time) { a.E.CreateTime = t }
func (a userInvitationAuditableEntity) SetUpdateTime(t time.Time) { a.E.UpdateTime = t }
func (a userInvitationAuditableEntity) GetVersion() int32         { return a.E.Version }
func (a userInvitationAuditableEntity) SetVersion(v int32)        { a.E.Version = v }
