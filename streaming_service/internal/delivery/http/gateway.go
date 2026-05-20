package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	grpcclient "github.com/Adiilkwz/music-grpc-go/streaming"
	"google.golang.org/grpc"
)

type Gateway struct {
	client grpcclient.StreamingServiceClient
}

func NewGateway(grpcAddr string) (*Gateway, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := grpcclient.NewStreamingServiceClient(conn)
	return &Gateway{client: client}, nil
}

func (g *Gateway) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stream/playlists", g.withCORS(g.handleGetPlaylists))
	mux.HandleFunc("/api/stream/likes", g.withCORS(g.handleGetLikes))

	srv := &http.Server{Addr: addr, Handler: mux}
	log.Printf("HTTP gateway listening on %s", addr)
	return srv.ListenAndServe()
}

func (g *Gateway) withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func (g *Gateway) handleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("playlist_id")
	if q == "" {
		http.Error(w, "missing param playlist_id", http.StatusBadRequest)
		return
	}
	resp, err := g.client.GetPlaylist(context.Background(), &grpcclient.GetPlaylistRequest{PlaylistId: parseInt64(q)})
	if err != nil {
		http.Error(w, fmt.Sprintf("gateway error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleGetLikes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("limit")
	limit := int32(20)
	if q != "" {
		limit = int32(parseInt(q, 20))
	}
	resp, err := g.client.GetLikedSongs(context.Background(), &grpcclient.GetLikedSongsRequest{Limit: limit, Offset: 0})
	if err != nil {
		http.Error(w, fmt.Sprintf("gateway error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func parseInt(s string, def int) int {
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return def
	}
	return v
}

func parseInt64(s string) int64 {
	var v int64
	_, _ = fmt.Sscanf(s, "%d", &v)
	return v
}
