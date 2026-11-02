package auth

func (s *Service) Reissue(userID int32, username string) (LoginResult, error) {
	return s.issue(userID, username)
}
