// TODO errutil
package errutil

import (
	"fmt"
)

type ErrReason int

const (
	ErrInvalidParameters = iota
	ErrSetupDir
	ErrOpenFile
	ErrNewFile
	ErrEmptyFile
	ErrInvalidComponent
	ErrSettingVersion
	ErrDownloadComponent
	ErrInstallComponent
	ErrUninstallComponent
	ErrUpdateComponent
)

// Error contains error reason and the actual error if any
// satisfies golang's error interface
type Error struct {
	reason ErrReason
	err    error
}

func New(reason ErrReason, err error) *Error {
	return &Error{reason: reason, err: err}
}

func (err Error) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s (%v)", err.description(), err.err)
	}
	return err.description()
}

func (err Error) description() string {
	switch err.reason {
	case ErrInvalidParameters:
		return "Invalid parameters"
	case ErrInvalidComponent:
		return "Invalid component"
	case ErrDownloadComponent:
		return "Failed to download component"
	case ErrInstallComponent:
		return "Failed to install component"
	case ErrUninstallComponent:
		return "Failed to uninstall component"
	case ErrOpenFile:
		return "Failed to open file"
	case ErrNewFile:
		return "Failed to create new file db"
	case ErrEmptyFile:
		return "Empty file"
	case ErrUpdateComponent:
		return "Failed to update component"
	case ErrSettingVersion:
		return "Error setting component's version"
	}

	return "Unknown error"
}
