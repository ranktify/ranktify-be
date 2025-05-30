package service

import (
	"context"
	"net/http"

	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
)

type RankingsService struct {
	RankingsDAO *dao.RankingsDao
	StreaksDAO  *dao.StreaksDAO
}

func NewRankingsService(rankingsDao *dao.RankingsDao, sDao *dao.StreaksDAO) *RankingsService {
	return &RankingsService{
		RankingsDAO: rankingsDao,
		StreaksDAO:  sDao,
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

func (s *RankingsService) RankSong(songID uint64, userID uint64, rank int) (int, content) {
	err := s.RankingsDAO.RankSong(songID, userID, rank)
	if err != nil {
		return http.StatusBadRequest, content{"error": "Failed to rank song"}
	}
	if err = s.StreaksDAO.RecordSongRank(context.Background(), userID); err != nil {
		return http.StatusBadRequest, content{"error": "Failed to record streak"}
	}
	return http.StatusOK, content{"Song ranked succesfully as a": rank}
}

func (s *RankingsService) DeleteRanking(rankingID uint64) (int, content) {
	err := s.RankingsDAO.DeleteRanking(rankingID)
	if err != nil {
		return http.StatusBadRequest, content{"error": "Failed to delete rank"}
	}
	return http.StatusOK, content{"Ranking deleted succesfully": rankingID}
}

func (s *RankingsService) UpdateRanking(rankingID uint64, rank int) (int, content) {
	err := s.RankingsDAO.UpdateRanking(rankingID, rank)
	if err != nil {
		return http.StatusBadRequest, content{"error": "Failed to update rank"}
	}
	return http.StatusOK, content{"Rank updated succesfully to": rank}

}

func (s *RankingsService) GetTopWeeklyRankedSongs(ctx context.Context) ([]model.Song, error) {
	songs, err := s.RankingsDAO.GetTopWeeklyRankedSongs(ctx)
	if err != nil {
		return nil, err
	}
	return songs, nil
}
