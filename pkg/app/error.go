package app

import (
	"errors"
	"fmt"
)

var (
	ErrExists        = errors.New("has already existed.")
	ErrNotExists     = errors.New("doesn't exist.")
	ErrInvalidParams = errors.New("contain invalid chars.")

	ErrUserExists    = fmt.Errorf("%w", ErrExists)
	ErrUserNotExists = fmt.Errorf("%w", ErrNotExists)

	ErrFolderExists    = fmt.Errorf("%w", ErrExists)
	ErrFolderNotExists = fmt.Errorf("%w", ErrNotExists)
	ErrListFolderEmpty = errors.New("doesn't have any folders.")

	ErrFileExists    = fmt.Errorf("%w", ErrExists)
	ErrFileNotExists = fmt.Errorf("%w", ErrNotExists)
	ErrListFileEmpty = errors.New("Warning: The folder is empty.")
)
