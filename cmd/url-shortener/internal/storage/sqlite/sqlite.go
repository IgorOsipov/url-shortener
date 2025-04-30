package sqlite

import (
	"fmt"
	"iosipoff/url-shortener/cmd/url-shortener/internal/storage"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type URL struct {
	ID    uint   `gorm:"primaryKey;column:id"`
	Alias string `gorm:"uniqueIndex;column:alias"`
	URL   string `gorm:"column:url"`
}

type Storage struct {
	db *gorm.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db, err := gorm.Open((sqlite.Open(storagePath)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.AutoMigrate(&URL{}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (uint, error) {
	const op = "storage.sqlite.SaveURL"

	url := URL{
		Alias: alias,
		URL:   urlToSave,
	}

	res := s.db.Create(&url)
	if res.Error != nil {
		return 0, fmt.Errorf("%s: %w", op, res.Error)
	}

	return url.ID, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	var url URL
	res := s.db.Where(&URL{Alias: alias}).First(&url)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, res.Error)
	}

	return url.URL, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	res := s.db.Delete(&URL{Alias: alias})
	if res.Error != nil {
		return fmt.Errorf("%s: %w", op, res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
