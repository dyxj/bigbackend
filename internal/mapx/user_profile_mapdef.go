package mapx

import (
	"github.com/dyxj/bigbackend/internal/sqlgen/user_profile_local/public/entity"
	"github.com/dyxj/bigbackend/internal/userprofile"
)

// goverter:converter
// goverter:output:file ./mapper/user_profile_mapper.go
// goverter:output:package mapper
// goverter:name UserProfile
// goverter:extend MapTime
// goverter:extend MapDate
// goverter:extend MapUUID
type UserProfileMapper interface {
	ModelToEntity(source userprofile.UserProfile) entity.UserProfile
}
