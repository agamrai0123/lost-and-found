package pkg

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func (s *ItemService) loadFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fileContent, err := os.Open(s.storePath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	data, err := io.ReadAll(fileContent)
	if err != nil {
		return err
	}

	var list []*Post
	if err := json.Unmarshal([]byte(data), &list); err != nil {
		return err
	}
	s.posts = make(map[string]*Post, len(list))
	for _, p := range list {
		// keep pointer in map
		s.posts[p.ID] = p
	}
	s.logger.Printf("loaded %d posts from %s", len(s.posts), s.storePath)
	return nil
}

func (s *ItemService) saveToFile() error {
	s.mu.RLock()
	list := make([]*Post, 0, len(s.posts))
	for _, p := range s.posts {
		list = append(list, p)
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.storePath)
	if dir == "." || dir == "" {
		dir = "."
	}
	tmpFile, err := os.CreateTemp(dir, "posts-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	tmpFile.Close()

	// atomic replace
	if err := os.Rename(tmpPath, s.storePath); err != nil {
		os.Remove(tmpPath)
		return err
	}
	s.logger.Printf("saved %d posts to %s", len(list), s.storePath)
	return nil
}
