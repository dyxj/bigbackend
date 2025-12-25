package app

import "github.com/dyxj/bigbackend/internal/userprofile"

func (s *Server) buildUserProfileHandlers() (
	*userprofile.CreatorHandler,
	*userprofile.GetterHandler,
) {
	mapper := &userprofile.UserProfileMapper{}
	repo := userprofile.NewCreatorSQLDB(s.logger)
	creator := userprofile.NewCreator(s.logger, repo, mapper)

	return userprofile.NewCreatorHandler(s.logger, s.dbConn, creator, mapper),
		userprofile.NewGetterHandler(s.logger)
}
