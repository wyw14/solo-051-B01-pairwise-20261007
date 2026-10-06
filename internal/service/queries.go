package service

import (
	"github.com/wyw14/cry-161/internal/model"
)

func (s *Service) ChamberSnapshot() model.ChamberSnapshot { return s.Chamber.Snapshot() }
func (s *Service) Events(limit int) []model.Event {
	values, err := s.Journal.Events(limit)
	if err != nil {
		return []model.Event{}
	}
	return values
}
