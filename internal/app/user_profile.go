package app

import (
	"github.com/dyxj/bigbackend/internal/user/profile"
)

func (s *Server) buildUserProfileHandlers() (
	*profile.CreatorHandler,
	*profile.GetterHandler,
) {
	mapper := &profile.UserProfileMapper{}

	cRepo := profile.NewCreatorSQLDB(s.logger)
	creator := profile.NewCreator(s.logger, cRepo, mapper)

	gRepo := profile.NewGetterSQLDB(s.logger, s.dbConn)
	getter := profile.NewGetter(s.logger, gRepo, mapper)

	return profile.NewCreatorHandler(s.logger, s.dbConn, creator, mapper),
		profile.NewGetterHandler(s.logger, getter, mapper)
}
