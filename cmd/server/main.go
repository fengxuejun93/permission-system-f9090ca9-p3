package main

import (
	"log"
	"net/http"
	"secondhand-trade/internal/handler"
	"secondhand-trade/internal/service"
	"strings"
)

func main() {
	itemService := service.NewItemService()
	itemHandler := handler.NewItemHandler(itemService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/items/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/items/")
		parts := strings.Split(path, "/")

		if len(parts) >= 2 {
			id := parts[0]
			action := parts[1]
			switch action {
			case "favorite":
				itemHandler.ToggleFavorite(w, r, id)
				return
			case "trade-intent":
				itemHandler.AddTradeIntent(w, r, id)
				return
			case "mark-communicated":
				itemHandler.MarkCommunicated(w, r, id)
				return
			case "offline":
				itemHandler.Offline(w, r, id)
				return
			case "relist":
				itemHandler.Relist(w, r, id)
				return
			}
		}

		if len(parts) == 1 && parts[0] != "" {
			id := parts[0]
			switch r.Method {
			case http.MethodGet:
				itemHandler.GetByID(w, r, id)
			case http.MethodPut:
				itemHandler.Update(w, r, id)
			case http.MethodDelete:
				itemHandler.Delete(w, r, id)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if path == "" || path == "/" {
			switch r.Method {
			case http.MethodGet:
				itemHandler.GetList(w, r)
			case http.MethodPost:
				itemHandler.Create(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/statistics", itemHandler.GetStatistics)
	mux.HandleFunc("/api/meta", itemHandler.GetMeta)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		itemHandler.ServeIndex(w, r)
	})

	log.Println("Server starting on :8080...")
	log.Println("Open http://localhost:8080 to view the application")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
