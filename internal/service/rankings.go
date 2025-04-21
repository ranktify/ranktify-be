package service

import (
	"net/http"

	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
)

type RankingsService struct {
	RankingsDAO *dao.RankingsDao
	FriendsDAO  *dao.FriendsDAO
}

func NewRankingsService(rankingsDao *dao.RankingsDao, friendsDao *dao.FriendsDAO) *RankingsService {
	return &RankingsService{
		RankingsDAO: rankingsDao,
		FriendsDAO:  friendsDao,
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
	friends, err := s.FriendsDAO.GetFriends(userID)
	if err != nil {
		return http.StatusNotFound, content{"error": "Failed to retrieve friends"}
	}
	var friendIDs []uint64
	for _, friend := range friends {
		friendIDs = append(friendIDs, uint64(friend.Id))
	}
	var rankedSongs []model.Rankings
	for _, friend := range friendIDs {
		rankedSong, err := s.RankingsDAO.GetRankedSongs(uint64(friend))
		if err != nil {
			return http.StatusNotFound, content{"error": "Failed to retrieve ranked songs"}
		}
		rankedSongs = append(rankedSongs, rankedSong...)
	}
	return http.StatusOK, content{"User's friends rankings": rankedSongs}
}
