package app

import (
	"errors"
	"fmt"
)

var (
	ErrExists        = errors.New("has already existed.")
	ErrNotExists     = errors.New("doesn't exist.")
	ErrInvalidParams = errors.New("contain invalid chars.")

	ErrFolderExists    = fmt.Errorf("%w", ErrExists)
	ErrFileExists      = fmt.Errorf("%w", ErrExists)
	ErrFolderNotExists = fmt.Errorf("%w", ErrNotExists)
	ErrFileNotExists   = fmt.Errorf("%w", ErrNotExists)
	ErrListFolderEmpty = errors.New("doesn't have any folders.")
	ErrListFileEmpty   = errors.New("The folder is empty.")
)
