package book

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"manimbook-reader/utils"
	"time"
)

// Book represents the structure of index.json
type Book struct {
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	PublishDate time.Time `json:"publishDate"`
	Chapters    []string  `json:"chapters"`
}

// ParseIndexJSON reads and parses index.json from the provided afero filesystem
func ParseIndexJSON(fs afero.Fs, book *Book) error {
	file, err := fs.Open("book/index.json")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&book); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func ValidateChapterFiles(fs afero.Fs, book *Book) error {
	for _, chapter := range book.Chapters {
		chapterPath := fmt.Sprintf("book/%s/index.html", chapter)
		exists, err := afero.Exists(fs, chapterPath)
		if err != nil {
			return fmt.Errorf("error checking file %s: %w", chapterPath, err)
		}
		if !exists {
			return fmt.Errorf("missing index.html in chapter folder: %s", chapterPath)
		}
	}
	return nil
}

func InitializeBook(path string, assets afero.Fs) (*Book, error) {
	var newBook Book
	err := utils.Unzip(path, assets, "book")
	if err != nil {
		return nil, fmt.Errorf("could not unzip file: %w", err)
	}
	err = ParseIndexJSON(assets, &newBook)
	if err != nil {
		return nil, fmt.Errorf("could not parse book: %w", err)
	}
	err = ValidateChapterFiles(assets, &newBook)
	if err != nil {
		return nil, fmt.Errorf("could not open book: %w", err)
	}
	return &newBook, nil
}
