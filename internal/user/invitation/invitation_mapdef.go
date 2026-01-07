package invitation

import "github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"

// goverter:converter
// goverter:output:file ./invitation_mapper.go
// goverter:name UserInvitationMapper
// goverter:extend github.com/dyxj/bigbackend/pkg/mapx:MapTime
// goverter:extend github.com/dyxj/bigbackend/pkg/mapx:MapUUID
type Mapper interface {
	// goverter:map StatusRaw Status
	ModelToEntity(source UserInvitation) entity.UserInvitation
	// goverter:map Status StatusRaw
	EntityToModel(source entity.UserInvitation) UserInvitation
	// goverter:ignoreMissing
	CreateRequestToModel(source CreateRequest) UserInvitation
}
