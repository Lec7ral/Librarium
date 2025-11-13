package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Lec7ral/fullAPI/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type cachingBookRepository struct {
	next  BookRepository
	cache *redis.Client
	ctx   context.Context
}

func NewCachingBookRepository(next BookRepository, cache *redis.Client) BookRepository {
	return &cachingBookRepository{
		next:  next,
		cache: cache,
		ctx:   context.Background(),
	}
}

func (r *cachingBookRepository) GetByID(id int64) (*models.Book, error) {
	key := fmt.Sprintf("book:%d", id)
	val, err := r.cache.Get(r.ctx, key).Result()
	if err == nil {
		var book models.Book
		if json.Unmarshal([]byte(val), &book) == nil {
			return &book, nil
		}
	}
	book, err := r.next.GetByID(id)
	if err != nil {
		return nil, err
	}
	jsonData, _ := json.Marshal(book)
	r.cache.Set(r.ctx, key, jsonData, 5*time.Minute)
	return book, nil
}
func (r *cachingBookRepository) Update(id int64, book models.Book) error {
	err := r.next.Update(id, book)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("book:%d", id)
	r.cache.Del(r.ctx, key)
	return nil
}

func (r *cachingBookRepository) Delete(id int64) error {
	err := r.next.Delete(id)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("book:%d", id)
	r.cache.Del(r.ctx, key)
	return nil
}

func (r *cachingBookRepository) Create(book models.Book) (int64, error) {
	return r.next.Create(book)
}
func (r *cachingBookRepository) Search(filter BookFilter, limit, offset int, sort, order string) ([]models.Book, int, error) {
	return r.next.Search(filter, limit, offset, sort, order)
}
