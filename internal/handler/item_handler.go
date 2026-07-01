package handler

import (
	"encoding/json"
	"net/http"
	"secondhand-trade/internal/model"
	"secondhand-trade/internal/service"
	"strconv"
)

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(s *service.ItemService) *ItemHandler {
	return &ItemHandler{service: s}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: "success",
		Data:    data,
	})
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
	})
}

func (h *ItemHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var query model.ItemQuery
	query.Keyword = r.URL.Query().Get("keyword")
	query.Category = r.URL.Query().Get("category")
	query.City = r.URL.Query().Get("city")
	query.Status = model.ItemStatus(r.URL.Query().Get("status"))

	items, err := h.service.GetList(query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.IncrementView(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.ItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" || req.Category == "" || req.City == "" || req.Condition == "" || req.Publisher == "" {
		writeError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	item, err := h.service.Create(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req model.ItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.service.Update(id, req)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, nil)
}

func (h *ItemHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats, err := h.service.GetStatistics()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *ItemHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.ToggleFavorite(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) AddTradeIntent(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.AddTradeIntent(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) MarkCommunicated(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.MarkCommunicated(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Offline(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.Offline(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Relist(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.Relist(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) GetMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	meta := map[string]interface{}{
		"categories": h.service.GetCategories(),
		"cities":     h.service.GetCities(),
		"conditions": h.service.GetConditions(),
		"statuses": []map[string]string{
			{"value": string(model.StatusActive), "label": "上架中"},
			{"value": string(model.StatusOffline), "label": "已下架"},
			{"value": string(model.StatusTraded), "label": "已置换"},
			{"value": string(model.StatusPending), "label": "待审核"},
		},
	}
	writeJSON(w, http.StatusOK, meta)
}

func (h *ItemHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func getIntQuery(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return result
}
