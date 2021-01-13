package dwn

import (
	"fmt"
	"path/filepath"
	"reflect"
)

var currentDwn DwnAttrs

// DwnAttrs list attributes needed for dwn
type DwnAttrs struct {
	DWN_URL      string
	DWN_FILENAME string
	DWN_BIN_NAME string

	DWN_FILENAME_CALC string
	DWN_FILENAME_BASE string
	DWN_FILENAME_EXT  string
}

// NewDwnAttrs enerates
func NewDwnAttrs() (DwnAttrs, error) {

	// 	DWN_URL:=getcourage.org	# Github URL to the file
	// DWN_FILENAME:=hello		# Github FileName
	// DWN_BIN_NAME:=?			# Local filename (the actual bin)

	// # calculated private variables
	// DWN_FILENAME_CALC=$(notdir $(DWN_URL)) # todo use this, so we dont need to pass in this anymore :)
	// DWN_FILENAME_BASE=$(shell basename -- $(DWN_FILENAME))
	// DWN_FILENAME_EXT := $(suffix $(DWN_FILENAME))
	// ifeq ($(DWN_FILENAME_EXT),)
	// 	DWN_FILENAME_EXT += NONE
	// endif

	DWN_URL := "getcourage.org"
	DWN_FILENAME := "hello"

	return DwnAttrs{
		DWN_URL:           DWN_URL,
		DWN_FILENAME:      DWN_FILENAME,
		DWN_BIN_NAME:      "",
		DWN_FILENAME_CALC: DWN_URL,
		DWN_FILENAME_BASE: DWN_FILENAME,
		DWN_FILENAME_EXT:  filepath.Ext(DWN_FILENAME),
	}, nil

}

func init() {
	var err error
	currentDwn, err = NewDwnAttrs()
	if err != nil {
		panic(err)
	}
}

// Print prints dwn attributes
func Print() error {

	dwnReflect := reflect.ValueOf(currentDwn)
	typeOfTarget := dwnReflect.Type()

	for i := 0; i < dwnReflect.NumField(); i++ {
		fmt.Printf("%s : %s \n", typeOfTarget.Field(i).Name, dwnReflect.Field(i).Interface())
	}

	return nil
}
