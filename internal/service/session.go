package service

import "github.com/bartosz-skejcik/go-analytics/internal/db/models"

func (s *Service) GetAllSessions() ([]models.Session, error) {
	s.db.Connect()
	defer s.db.Db.Close()

	query := `
  select id, anonymous_id, timestamp, referrer, screen_width, ip, user_agent, country, country_code, os, browser from session;
  `
	rows, err := s.db.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		err := rows.Scan(
			&s.Id, &s.AnonymousID, &s.Timestamp, &s.Referrer, &s.ScreenWidth, &s.IP, &s.UserAgent, &s.Country, &s.CountryCode, &s.OS, &s.Browser,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}
