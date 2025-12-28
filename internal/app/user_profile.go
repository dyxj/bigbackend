package app

import (
	"github.com/dyxj/bigbackend/internal/userprofile"
)

func (s *Server) buildUserProfileHandlers() (
	*userprofile.CreatorHandler,
	*userprofile.GetterHandler,
) {
	mapper := &userprofile.UserProfileMapper{}

	cRepo := userprofile.NewCreatorSQLDB(s.logger)
	creator := userprofile.NewCreator(s.logger, cRepo, mapper)

	gRepo := userprofile.NewGetterSQLDB(s.logger, s.dbConn)
	getter := userprofile.NewGetter(s.logger, gRepo, mapper)

	return userprofile.NewCreatorHandler(s.logger, s.dbConn, creator, mapper),
		userprofile.NewGetterHandler(s.logger, getter, mapper)
}
