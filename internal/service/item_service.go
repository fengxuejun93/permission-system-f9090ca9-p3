package service

import (
	"errors"
	"fmt"
	"secondhand-trade/internal/data"
	"secondhand-trade/internal/model"
	"strings"
	"sync"
	"time"
)

type ItemService struct {
	items map[string]*model.Item
	mu    sync.RWMutex
}

func NewItemService() *ItemService {
	s := &ItemService{
		items: make(map[string]*model.Item),
	}
	for _, item := range data.MockItems {
		s.items[item.ID] = item
	}
	return s
}

func (s *ItemService) GetList(query model.ItemQuery) ([]*model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Item
	for _, item := range s.items {
		if query.Keyword != "" {
			keyword := strings.ToLower(query.Keyword)
			titleMatch := strings.Contains(strings.ToLower(item.Title), keyword)
			descMatch := strings.Contains(strings.ToLower(item.Description), keyword)
			publisherMatch := strings.Contains(strings.ToLower(item.Publisher), keyword)
			if !titleMatch && !descMatch && !publisherMatch {
				continue
			}
		}
		if query.Category != "" && item.Category != query.Category {
			continue
		}
		if query.City != "" && item.City != query.City {
			continue
		}
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *ItemService) GetByID(id string) (*model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (s *ItemService) Create(req model.ItemCreateRequest) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	item := &model.Item{
		ID:              id,
		Title:           req.Title,
		Category:        req.Category,
		City:            req.City,
		Condition:       req.Condition,
		Publisher:       req.Publisher,
		DesiredCategory: req.DesiredCategory,
		Description:     req.Description,
		ViewCount:       0,
		FavoriteCount:   0,
		TradeIntentCount: 0,
		Status:          model.StatusActive,
		IsFavorited:     false,
		HasCommunicated: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	s.items[id] = item
	return item, nil
}

func (s *ItemService) Update(id string, req model.ItemUpdateRequest) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.Title = req.Title
	item.Category = req.Category
	item.City = req.City
	item.Condition = req.Condition
	item.DesiredCategory = req.DesiredCategory
	item.Description = req.Description
	item.UpdatedAt = time.Now()

	return item, nil
}

func (s *ItemService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.items[id]
	if !exists {
		return errors.New("item not found")
	}
	delete(s.items, id)
	return nil
}

func (s *ItemService) GetStatistics() (*model.Statistics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &model.Statistics{}
	for _, item := range s.items {
		stats.TotalItems++
		stats.TotalViews += item.ViewCount
		stats.TotalFavorites += item.FavoriteCount
		stats.TotalIntents += item.TradeIntentCount

		switch item.Status {
		case model.StatusActive:
			stats.ActiveItems++
		case model.StatusOffline:
			stats.OfflineItems++
		case model.StatusTraded:
			stats.TradedItems++
		}
	}
	return stats, nil
}

func (s *ItemService) IncrementView(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return errors.New("item not found")
	}
	item.ViewCount++
	item.UpdatedAt = time.Now()
	return nil
}

func (s *ItemService) ToggleFavorite(id string) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	if item.IsFavorited {
		item.FavoriteCount--
		item.IsFavorited = false
	} else {
		item.FavoriteCount++
		item.IsFavorited = true
	}
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *ItemService) AddTradeIntent(id string) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	if item.Status == model.StatusOffline {
		return nil, errors.New("已下架的货品不能发起置换意向")
	}

	item.TradeIntentCount++
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *ItemService) MarkCommunicated(id string) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.HasCommunicated = true
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *ItemService) Offline(id string) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.Status = model.StatusOffline
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *ItemService) Relist(id string) (*model.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	item.Status = model.StatusActive
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *ItemService) GetCategories() []string {
	return data.Categories
}

func (s *ItemService) GetCities() []string {
	return data.Cities
}

func (s *ItemService) GetConditions() []string {
	return data.Conditions
}
