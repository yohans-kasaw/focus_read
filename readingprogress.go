package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Progress struct {
	Index int
	Path  string
}

type ProgressStore struct {
	Progresses []Progress
	Index      int
}

func NewProgressStore() (*ProgressStore, error) {
	m := &ProgressStore{
		Progresses: []Progress{},
		Index:      0,
	}

	file, err := os.ReadFile(filepath.Join(cacheDir, progressFile))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error: reading progress store %v\n", err)
	}

	if err == nil {
		if err := json.Unmarshal(file, &m); err != nil {
			return nil, fmt.Errorf("Error Unmarshaling progress store%v\n", err)
		}
	}

	return m, nil
}

func (m *ProgressStore) GetProgress(path string) *Progress{
	m.SetCurrent(path)
	return &m.Progresses[m.Index]
}

func (m *ProgressStore) SetCurrent(path string) {
	for i := range m.Progresses {
		if m.Progresses[i].Path == path {
			m.Index = i
			return
		}
	}

	m.Progresses = append(m.Progresses, Progress{Path: path, Index: 0})
	m.Index = len(m.Progresses) - 1
}

func (m *ProgressStore) SaveProgress() error {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("Error creating cache folder: %v\n", err)
	}

	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Error marshaling Progress Store %v\n", err)
	}

	err = os.WriteFile(filepath.Join(cacheDir, progressFile), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write Progress Store to file: %w", err)
	}
	return nil
}
