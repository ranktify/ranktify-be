package service

import (
	"context"
	"net/http"

	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
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

func (s *RankingsService) GetFriendsRankedSongsWithNoUserRank(userID uint64) (int, content) {
	rankings, err := s.RankingsDAO.GetFriendsRankedSongsWithNoUserRank(userID)
	if err != nil {
		return http.StatusNotFound, content{"error": "Failed to retrieve rankings"}
	}
	return http.StatusOK, content{"User's friends songs": rankings}
}

func (s *RankingsService) GetTopWeeklyRankedSongs(ctx context.Context) ([]model.Song, error) {
	songs, err := s.RankingsDAO.GetTopWeeklyRankedSongs(ctx)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
