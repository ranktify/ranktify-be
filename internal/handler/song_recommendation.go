package handler

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
	"github.com/ranktify/ranktify-be/internal/spotify"
)

type SongRecommendationHandler struct {
	RankingsDAO *dao.RankingsDao
}

func NewSongRecommendationHandler(rankingsDAO *dao.RankingsDao) *SongRecommendationHandler {
	return &SongRecommendationHandler{
		RankingsDAO: rankingsDAO,
	}
}

func (h *SongRecommendationHandler) SongRecommendation(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	rawToken, ok := c.Get("spotifyToken")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
		return
	}
	accessToken := rawToken.(string)

	randomSongs, err := spotify.GetRandomSongs(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	randomGenre := spotify.GetRandomGenre()

	randomSongsGenre, err := spotify.GetRandomSongsByGenre(c.Request.Context(), accessToken, randomGenre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	friendsSongs, err := h.RankingsDAO.GetFriendsRankedSongsWithNoUserRank(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	songList := append(*randomSongs, *randomSongsGenre...)

	// Shuffle in-place
	var recommendedSongs []model.Song
	for _, song := range songList {
		isRanked, err := h.RankingsDAO.CheckIfSongIsRanked(song.SpotifyID, userID)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"error": err.Error(),
			})
			continue
		}

		if !isRanked {
			recommendedSongs = append(recommendedSongs, song)
		}
		if len(recommendedSongs) >= 15 {
			break
		}
	}
	rand.Shuffle(len(songList), func(i, j int) {
		songList[i], songList[j] = songList[j], songList[i]
	})

	rand.Shuffle(len(friendsSongs), func(i, j int) {
		friendsSongs[i], friendsSongs[j] = friendsSongs[j], friendsSongs[i]
	})

	if len(friendsSongs) >= 5 {
		friendsSongs = friendsSongs[0:5]
		recommendedSongs = recommendedSongs[0:5]
	} else {
		limit := 10 - len(friendsSongs)
		if limit > len(recommendedSongs) {
			limit = len(recommendedSongs)
		}
		recommendedSongs = recommendedSongs[:limit]
		// recommendedSongs = recommendedSongs[0:(10 - len(friendsSongs))]
	}

	var chosenSongs []model.Song
	for _, song := range recommendedSongs {
		h.RankingsDAO.StoreSongInDB(
			song.SpotifyID,
			song.Title,
			song.Artist,
			song.Album,
			song.ReleaseDate,
			song.Genre,
			song.CoverURI,
			song.PreviewURI,
		)
		songWithID, _ := h.RankingsDAO.GetSongBySpotifyID(song.SpotifyID)
		chosenSongs = append(chosenSongs, songWithID)

	}

	c.JSON(http.StatusOK, gin.H{
		"Recommended Songs":  chosenSongs,
		"Songs From Friends": friendsSongs,
	})

}
