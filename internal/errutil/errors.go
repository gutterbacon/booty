// TODO errutil
package errutil

import (
	"fmt"
)

type ErrReason int

const (
	ErrInvalidParameters = iota
	ErrSetupDir
	Err
)

// Error contains error reason and the actual error if any
// satisfies golang's error interface
type Error struct {
	Reason ErrReason
	Err    error
}

func (err Error) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s (%v)", err.description(), err.Err)
	}
	return err.description()
}


func (err Error) description() string {
	switch err.Reason {
	case ErrInvalidParameters:
		return "Invalid parameters"
	}

	return "Unknown error"
}
