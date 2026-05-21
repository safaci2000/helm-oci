package bookmark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

type Bookmark struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type fileData struct {
	Bookmarks []Bookmark `yaml:"bookmarks"`
}

type Store struct {
	path string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) load() (fileData, error) {
	var data fileData

	raw, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return data, nil
	}
	if err != nil {
		return data, fmt.Errorf("reading bookmarks: %w", err)
	}

	if err := yaml.Unmarshal(raw, &data); err != nil {
		return data, fmt.Errorf("parsing bookmarks: %w", err)
	}
	return data, nil
}

func (s *Store) save(data fileData) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("creating bookmark directory: %w", err)
	}

	raw, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling bookmarks: %w", err)
	}
	return os.WriteFile(s.path, raw, 0o644)
}

func (s *Store) Add(name, url string) error {
	if !strings.HasPrefix(url, "oci://") {
		return fmt.Errorf("URL must start with oci://, got %q", url)
	}

	data, err := s.load()
	if err != nil {
		return err
	}

	for _, b := range data.Bookmarks {
		if b.Name == name {
			return fmt.Errorf("bookmark %q already exists", name)
		}
	}

	data.Bookmarks = append(data.Bookmarks, Bookmark{Name: name, URL: url})
	return s.save(data)
}

func (s *Store) Remove(name string) error {
	data, err := s.load()
	if err != nil {
		return err
	}

	for i, b := range data.Bookmarks {
		if b.Name == name {
			data.Bookmarks = append(data.Bookmarks[:i], data.Bookmarks[i+1:]...)
			return s.save(data)
		}
	}
	return fmt.Errorf("bookmark %q not found", name)
}

func (s *Store) Get(name string) (Bookmark, error) {
	data, err := s.load()
	if err != nil {
		return Bookmark{}, err
	}

	for _, b := range data.Bookmarks {
		if b.Name == name {
			return b, nil
		}
	}
	return Bookmark{}, fmt.Errorf("bookmark %q not found", name)
}

func (s *Store) List() ([]Bookmark, error) {
	data, err := s.load()
	if err != nil {
		return nil, err
	}
	return data.Bookmarks, nil
}
