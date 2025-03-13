package service

import (
	"fmt"
	"os"

	"github.com/fatalistix/terminal-file-explorer/internal/model"
)

type DirectoryLoadService struct{}

func (d *DirectoryLoadService) LoadDirectory(dir string) (model.Directory, error) {
	const op = "service.DirectoryLoadService.LoadDirectory"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return model.Directory{}, fmt.Errorf("%s: unable to read directory: %w", op, err)
	}

	entriesStr := make([]string, 0, len(entries))
	for _, v := range entries {
		entriesStr = append(entriesStr, v.Name())
	}

	return model.Directory{
		Entries: entriesStr,
	}, nil
}
