package service

import "github.com/bartosz-skejcik/go-analytics/internal/db/models"

func (s *Service) GetAllSessions() ([]models.Session, error) {
	s.db.Connect()
	pDb, err := s.db.Client.DB()
	if err != nil {
		panic(err)
	}

	defer pDb.Close()

	var sessions []models.Session
	s.db.Client.Model(&models.Session{})
	results := s.db.Client.Find(&sessions)

	return sessions, results.Error
}

func (s *Service) SessionCreateDtoToModel(dto *models.SessionCreateDto) (*models.Session, error) {
	// create the session in the database (get the session id back)
	// go through each ...

	return nil, nil
}
