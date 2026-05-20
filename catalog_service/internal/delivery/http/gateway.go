package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	grpcclient "github.com/Adiilkwz/music-grpc-go/catalog"
	"google.golang.org/grpc"
)

type Gateway struct {
	grpcAddr string
	client   grpcclient.CatalogServiceClient
}

func NewGateway(grpcAddr string) (*Gateway, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := grpcclient.NewCatalogServiceClient(conn)
	return &Gateway{grpcAddr: grpcAddr, client: client}, nil
}

func (g *Gateway) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", g.withCORS(g.handleSearch))

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
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

func (g *Gateway) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "missing query param 'q'", http.StatusBadRequest)
		return
	}
	limit := int32(10)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := g.client.SearchCatalog(ctx, &grpcclient.SearchCatalogRequest{Query: q, Limit: limit})
	if err != nil {
		http.Error(w, fmt.Sprintf("gateway error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
