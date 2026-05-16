package http

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"api_gateway/internal/middleware"

	"github.com/Adiilkwz/music-grpc-go/catalog"
	"github.com/Adiilkwz/music-grpc-go/streaming"
)

type CreatePlaylistInput struct {
	Title string `json:"title" binding:"required"`
}

type ModifyPlaylistInput struct {
	SongID int64 `json:"song_id" binding:"required"`
}

type LikeSongInput struct {
	SongID int64 `json:"song_id" binding:"required"`
}

type RecordPlayInput struct {
	SongID int64 `json:"song_id" binding:"required"`
}

func RegisterStreamRoutes(rg *gin.RouterGroup, streamClient streaming.StreamingServiceClient, catalogClient catalog.CatalogServiceClient) {
	playlistGroup := rg.Group("/playlists")
	{
		playlistGroup.POST("/", createPlaylistHandler(streamClient))
		playlistGroup.GET("/:id", getPlaylistHandler(streamClient))
		playlistGroup.POST("/:id/songs", addSongToPlaylistHandler(streamClient))
		playlistGroup.DELETE("/:id/songs/:song_id", removeSongFromPlaylistHandler(streamClient))
		playlistGroup.DELETE("/:id", deletePlaylistHandler(streamClient))
	}

	likesGroup := rg.Group("/likes")
	{
		likesGroup.POST("/", likeSongHandler(streamClient))
		likesGroup.DELETE("/:song_id", unlikeSongHandler(streamClient))
		likesGroup.GET("/", getLikedSongsEnrichedHandler(streamClient, catalogClient))
	}

	audioGroup := rg.Group("/audio")
	{
		audioGroup.GET("/stream/:song_id", streamAudioHandler(streamClient))
		audioGroup.POST("/history", recordPlayHandler(streamClient))
		audioGroup.GET("/history", getUserHistoryHandler(streamClient))
		audioGroup.GET("/trending", getTrendingHandler(streamClient))
	}
}

func createPlaylistHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreatePlaylistInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.CreatePlaylist(ctx, &streaming.CreatePlaylistRequest{Title: input.Title})
		handleProtectedStreamResponse(c, resp, err, http.StatusCreated)
	}
}

func getPlaylistHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.GetPlaylist(ctx, &streaming.GetPlaylistRequest{PlaylistId: id})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func addSongToPlaylistHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		var input ModifyPlaylistInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.AddSongToPlaylist(ctx, &streaming.ModifyPlaylistRequest{
			PlaylistId: playlistID,
			SongId:     input.SongID,
		})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func removeSongFromPlaylistHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		songID, err2 := strconv.ParseInt(c.Param("song_id"), 10, 64)
		if err != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IDs"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.RemoveSongFromPlaylist(ctx, &streaming.ModifyPlaylistRequest{
			PlaylistId: playlistID,
			SongId:     songID,
		})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func deletePlaylistHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.DeletePlaylist(ctx, &streaming.DeletePlaylistRequest{PlaylistId: id})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func likeSongHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LikeSongInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.LikeSong(ctx, &streaming.LikeSongRequest{SongId: input.SongID})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func unlikeSongHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID, err := strconv.ParseInt(c.Param("song_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.UnlikeSong(ctx, &streaming.LikeSongRequest{SongId: songID})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func getLikedSongsEnrichedHandler(streamClient streaming.StreamingServiceClient, catalogClient catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
		offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)

		likesResp, err := streamClient.GetLikedSongs(ctx, &streaming.GetLikedSongsRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		var enrichedSongs []interface{}
		for _, songID := range likesResp.SongIds {
			songDetails, err := catalogClient.GetSong(c.Request.Context(), &catalog.GetSongRequest{Id: songID})
			if err == nil {
				enrichedSongs = append(enrichedSongs, songDetails)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"total_likes": len(likesResp.SongIds),
			"songs":       enrichedSongs,
		})
	}
}

func recordPlayHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RecordPlayInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.RecordPlay(ctx, &streaming.RecordPlayRequest{SongId: input.SongID})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func getUserHistoryHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := middleware.GetGrpcContext(c)
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

		resp, err := client.GetUserHistory(ctx, &streaming.GetUserHistoryRequest{Limit: int32(limit)})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func getTrendingHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
		resp, err := client.GetTrendingSongs(c.Request.Context(), &streaming.GetTrendingRequest{Limit: int32(limit)})
		handleProtectedStreamResponse(c, resp, err, http.StatusOK)
	}
}

func streamAudioHandler(client streaming.StreamingServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		songID, err := strconv.ParseInt(c.Param("song_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		stream, err := client.StreamAudio(c.Request.Context(), &streaming.StreamRequest{SongId: songID})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.Writer.Header().Set("Content-Type", "audio/mpeg")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error receiving audio stream from gRPC: %v", err)
				break
			}

			c.Writer.Write(resp.GetAudioChunk())
			c.Writer.Flush()
		}
	}
}

func handleProtectedStreamResponse(c *gin.Context, resp interface{}, err error, successStatus int) {
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Action requires appropriate privileges"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(successStatus, resp)
}
