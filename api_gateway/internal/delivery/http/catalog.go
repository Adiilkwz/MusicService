package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"api_gateway/internal/middleware"

	"github.com/Adiilkwz/music-grpc-go/catalog"
)

type CreateArtistInput struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio"`
}

type CreateAlbumInput struct {
	ArtistID    int64  `json:"artist_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	ReleaseYear int32  `json:"release_year" binding:"required"`
}

type CreateSongInput struct {
	AlbumID  int64  `json:"album_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Duration int32  `json:"duration" binding:"required"`
	Genre    string `json:"genre" binding:"required"`
}

type UpdateSongInput struct {
	Title string `json:"title" binding:"required"`
	Genre string `json:"genre" binding:"required"`
}

func RegisterPublicCatalogRoutes(rg *gin.RouterGroup, client catalog.CatalogServiceClient) {
	catalogGroup := rg.Group("/catalog")
	{
		catalogGroup.GET("/search", searchCatalogHandler(client))
		catalogGroup.GET("/artists/:id", getArtistHandler(client))
		catalogGroup.GET("/artists/:id/albums", getAlbumsByArtistHandler(client))
		catalogGroup.GET("/albums/:id", getAlbumHandler(client))
		catalogGroup.GET("/albums/:id/songs", getSongsByAlbumHandler(client))
		catalogGroup.GET("/songs/:id", getSongHandler(client))
		catalogGroup.GET("/genres/:genre", getSongsByGenreHandler(client))
	}
}

func RegisterProtectedCatalogRoutes(rg *gin.RouterGroup, client catalog.CatalogServiceClient) {
	catalogGroup := rg.Group("/catalog")
	{
		catalogGroup.POST("/artists", createArtistHandler(client))
		catalogGroup.POST("/albums", createAlbumHandler(client))
		catalogGroup.POST("/songs", createSongHandler(client))
		catalogGroup.PUT("/songs/:id", updateSongHandler(client))
		catalogGroup.DELETE("/songs/:id", deleteSongHandler(client))
	}
}

func searchCatalogHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query 'q' is required"})
			return
		}

		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
		resp, err := client.SearchCatalog(c.Request.Context(), &catalog.SearchCatalogRequest{
			Query: query,
			Limit: int32(limit),
		})
		if err != nil {
			st, _ := status.FromError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"artists": resp.Artists, "albums": resp.Albums, "songs": resp.Songs})
	}
}

func getArtistHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
		}

		resp, err := client.GetArtist(c.Request.Context(), &catalog.GetArtistRequest{Id: id})
		handleGrpcResponse(c, resp, err)
	}
}

func getAlbumsByArtistHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		artistID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
		}

		resp, err := client.GetAlbumsByArtist(c.Request.Context(), &catalog.GetAlbumsByArtistRequest{ArtistId: artistID})
		handleGrpcResponse(c, resp, err)
	}
}

func getAlbumHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}

		resp, err := client.GetAlbum(c.Request.Context(), &catalog.GetAlbumRequest{Id: id})
		handleGrpcResponse(c, resp, err)
	}
}

func getSongsByAlbumHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		albumID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}

		resp, err := client.GetSongsByAlbum(c.Request.Context(), &catalog.GetSongsByAlbumRequest{AlbumId: albumID})
		handleGrpcResponse(c, resp, err)
	}
}

func getSongHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		resp, err := client.GetSong(c.Request.Context(), &catalog.GetSongRequest{Id: id})
		handleGrpcResponse(c, resp, err)
	}
}

func getSongsByGenreHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		genre := c.Param("genre")
		limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)

		resp, err := client.GetSongsByGenre(c.Request.Context(), &catalog.GetSongsByGenreRequest{
			Genre: genre,
			Limit: int32(limit),
		})
		handleGrpcResponse(c, resp, err)
	}
}

func createArtistHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateArtistInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.CreateArtist(ctx, &catalog.CreateArtistRequest{
			Name: input.Name,
			Bio:  input.Bio,
		})
		handleProtectedGrpcResponse(c, resp, err, http.StatusCreated)
	}
}

func createAlbumHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateAlbumInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.CreateAlbum(ctx, &catalog.CreateAlbumRequest{
			ArtistId:    input.ArtistID,
			Title:       input.Title,
			ReleaseYear: input.ReleaseYear,
		})
		handleProtectedGrpcResponse(c, resp, err, http.StatusCreated)
	}
}

func createSongHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateSongInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.CreateSong(ctx, &catalog.CreateSongRequest{
			AlbumId:  input.AlbumID,
			Title:    input.Title,
			Duration: input.Duration,
			Genre:    input.Genre,
		})
		handleProtectedGrpcResponse(c, resp, err, http.StatusCreated)
	}
}

func updateSongHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		var input UpdateSongInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.UpdateSong(ctx, &catalog.UpdateSongRequest{
			Id:    id,
			Title: input.Title,
			Genre: input.Genre,
		})
		handleProtectedGrpcResponse(c, resp, err, http.StatusOK)
	}
}

func deleteSongHandler(client catalog.CatalogServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
			return
		}

		ctx := middleware.GetGrpcContext(c)
		resp, err := client.DeleteSong(ctx, &catalog.DeleteSongRequest{Id: id})
		handleProtectedGrpcResponse(c, resp, err, http.StatusOK)
	}
}

func handleGrpcResponse(c *gin.Context, resp interface{}, err error) {
	if err != nil {
		st, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func handleProtectedGrpcResponse(c *gin.Context, resp interface{}, err error, successStatus int) {
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.PermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": "Action requires artist or admin privileges"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}
	c.JSON(successStatus, resp)
}
