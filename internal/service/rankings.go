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