package userprofile

import (
	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
)

// goverter:converter
// goverter:output:file ./userprofile_mapper.go
// goverter:name UserProfileMapper
// goverter:extend github.com/dyxj/bigbackend/pkg/mapx:MapTime
// goverter:extend github.com/dyxj/bigbackend/pkg/mapx:MapDate
// goverter:extend github.com/dyxj/bigbackend/pkg/mapx:MapUUID
type Mapper interface {
	ModelToEntity(source UserProfile) entity.UserProfile
	EntityToModel(source entity.UserProfile) UserProfile
}
