package service

import (
	"net/http"

	"github.com/ranktify/ranktify-be/internal/dao"
)

type RankingsService struct {
	RankingsDAO *dao.RankingsDao
}

func NewRankingsService(rankingsDao *dao.RankingsDao) *RankingsService {
	return &RankingsService{
		RankingsDAO: rankingsDao,
	}
}

func (s *RankingsService) GetRankedSongs(userID uint64) (int, content) {
	rankings, err := s.RankingsDAO.GetRankedSongs(userID)
	if err != nil {
		return http.StatusNotFound, content{"error": "Failed to retrieve rankings"}
	}
	return http.StatusOK, content{"rankings": rankings}
}

func (s *RankingsService) GetFriendsRankedSongs(userID uint64) (int, content) {
	rankings, err := s.RankingsDAO.GetFriendsRankedSongs(userID)
	if err != nil {
		return http.StatusNotFound, content{"error": "Failed to retrieve rankings"}
	}
	return http.StatusOK, content{"User's friends rankings": rankings}
}
