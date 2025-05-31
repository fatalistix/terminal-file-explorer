package service

import (
	"errors"
	"fmt"
	"github.com/fatalistix/terminal-file-explorer/internal/model"
	"os"
	"path/filepath"
	"strings"
)

type StateHolder struct {
	state model.State
}

func NewStateHolder(startPath string) (*StateHolder, error) {
	const op = "service.NewStateHolder"

	currentPath, err := normalizePath(startPath)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to normalize path %s: %w", op, startPath, err)
	}

	currentDirEntries, err := readDir(currentPath)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to read directory at %s: %w", op, currentPath, err)
	}

	currentViewContent := &model.DirectoryContent{
		Directory: model.Directory{
			Entries: currentDirEntries,
		},
		SelectedId: 0,
	}

	var previousViewContent model.Content
	if hasParent(currentPath) {
		previousPath := getParent(currentPath)
		previousDirEntries, err := readDir(previousPath)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to read directory at %s: %w", op, previousPath, err)
		}

		currentDirName := getFilename(currentPath)
		previousDirSelectedId := -1
		for i, e := range previousDirEntries {
			if e.Name == currentDirName {
				previousDirSelectedId = i
			}
		}
		if previousDirSelectedId == -1 {
			panic("Unable to find previous directory entry")
		}

		previousViewContent = &model.DirectoryContent{
			Directory: model.Directory{
				Entries: previousDirEntries,
			},
			SelectedId: previousDirSelectedId,
		}
	} else {
		previousViewContent = &model.TextContent{
			Text: "No files under ROOT",
		}
	}

	nextViewContent := &model.EmptyContent{}

	state := model.State{
		PrevViewContent: previousViewContent,
		CurrViewContent: currentViewContent,
		NextViewContent: nextViewContent,
		CurrentPath:     currentPath,
	}

	return &StateHolder{state: state}, nil
}

func normalizePath(path string) (string, error) {
	const op = "service.normalizePath"

	normalizedPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("%s: unable to convert path %s to absolute: %w", op, path, err)
	}

	if strings.HasSuffix(normalizedPath, string(os.PathSeparator)) {
		normalizedPath = strings.TrimSuffix(normalizedPath, string(os.PathSeparator))
	}

	return normalizedPath, nil
}

func readDir(path string) ([]model.DirectoryEntry, error) {
	const op = "service.loadCurrentDir"

	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%s: file %s does not exists: %w", op, path, err)
		} else {
			return nil, fmt.Errorf("%s: error searching file %s: %w", op, path, err)
		}
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s: file %s is not a directory", op, path)
	}

	if strings.HasSuffix(path, string(os.PathSeparator)) {
		path = strings.TrimSuffix(path, string(os.PathSeparator))
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("%s: error reading directory %s: %w", op, path, err)
	}

	return mapOsEntriesToModelEntries(entries), nil
}

func hasParent(path string) bool {
	return filepath.Dir(path) != path
}

func getParent(path string) string {
	return filepath.Dir(path)
}

func getFilename(path string) string {
	return filepath.Base(path)
}

func mapOsEntryToModelEntry(osEntry os.DirEntry) model.DirectoryEntry {
	return model.DirectoryEntry{
		Name: osEntry.Name(),
	}
}

func mapOsEntriesToModelEntries(osEntries []os.DirEntry) []model.DirectoryEntry {
	entries := make([]model.DirectoryEntry, len(osEntries))
	for i, entry := range osEntries {
		entries[i] = mapOsEntryToModelEntry(entry)
	}
	return entries
}

func findIdByName(entries []model.DirectoryEntry, name string) int {
	for i, entry := range entries {
		if entry.Name == name {
			return i
		}
	}

	return -1
}

func (h *StateHolder) GetState() model.State {
	return h.state
}

func (h *StateHolder) MoveDown() {
	s, ok := h.state.CurrViewContent.(*model.DirectoryContent)
	if !ok {
		panic("Unexpected content type")
	}

	s.SelectedId = min(s.SelectedId+1, len(s.Directory.Entries)-1)
}

func (h *StateHolder) MoveUp() {
	s, ok := h.state.CurrViewContent.(*model.DirectoryContent)
	if !ok {
		panic("Unexpected content type")
	}

	s.SelectedId = max(s.SelectedId-1, 0)
}

func (h *StateHolder) Select() {
	s, ok := h.state.CurrViewContent.(*model.DirectoryContent)
	if !ok {
		panic("Unexpected content type")
	}

	h.state.PrevViewContent = s
	h.state.CurrentPath += "/" + s.Directory.Entries[s.SelectedId].Name

	currDirEntries, err := readDir(h.state.CurrentPath)
	if err != nil {
		panic(err)
	}

	content := &model.DirectoryContent{
		Directory: model.Directory{
			Entries: currDirEntries,
		},
		SelectedId: 0,
	}

	h.state.CurrViewContent = content
}

func (h *StateHolder) GoToPreviousDirectory() {
	if !hasParent(h.state.CurrentPath) {
		return
	}

	h.state.CurrViewContent = h.state.PrevViewContent
	h.state.CurrentPath = getParent(h.state.CurrentPath)

	prevDirPath := getParent(h.state.CurrentPath)
	prevDirEntries, err := readDir(prevDirPath)
	if err != nil {
		panic(err)
	}

	prevDirName := getFilename(h.state.CurrentPath)
	prevDirSelectedId := findIdByName(prevDirEntries, prevDirName)
	if prevDirSelectedId == -1 {
		panic("Unable to find previous directory entry")
	}

	content := &model.DirectoryContent{
		Directory: model.Directory{
			Entries: prevDirEntries,
		},
		SelectedId: prevDirSelectedId,
	}

	h.state.PrevViewContent = content
}
