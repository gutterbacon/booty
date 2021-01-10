package gitr

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-git/go-git/v5"
)

// Here lies gitr definitions.

var currentGitr GitrAttrs

// GitrAttrs represents the attributes needed for gitr
type GitrAttrs struct {
	GITR_SERVER                string
	GITR_ORG_UPSTREAM          string
	GITR_ORG_FORK              string
	GITR_USER                  string
	GITR_REPO_NAME             string
	GITR_REPO_UPSTREAM_ABS_URL string
	GITR_REPO_ABS_URL          string
	GITR_REPO_ABS_FSPATH       string
	GITR_BRANCH_NAME           string
	GITR_COMMIT_MESSAGE        string
	repoController             *git.Repository
}

// NewGitrAttrs generates gitr attributes
func NewGitrAttrs() (GitrAttrs, error) {
	wd, _ := os.Getwd()
	repo, err := git.PlainOpen(wd)
	if err != nil {
		return GitrAttrs{}, err
	}
	// abs, _ := filepath.Abs(dir)

	GITR_SERVER := "github.com"
	GITR_ORG_UPSTREAM := "getcouragenow"
	GITR_ORG_FORK := filepath.Base(filepath.Dir(wd)) //$(shell basename $(dir $(abspath $(dir $$PWD))))
	GITR_USER := filepath.Base(filepath.Dir(wd))
	GITR_REPO_NAME := filepath.Base(wd)

	return GitrAttrs{
		GITR_SERVER:                GITR_SERVER,
		GITR_ORG_UPSTREAM:          GITR_ORG_UPSTREAM,
		GITR_ORG_FORK:              GITR_ORG_FORK, //$(shell basename $(dir $(abspath $(dir $$PWD))))
		GITR_USER:                  GITR_USER,
		GITR_REPO_NAME:             GITR_REPO_NAME,
		GITR_REPO_UPSTREAM_ABS_URL: fmt.Sprintf("https://%s/%s/%s", GITR_SERVER, GITR_ORG_UPSTREAM, GITR_REPO_NAME),
		GITR_REPO_ABS_URL:          fmt.Sprintf("https://%s/%s/%s", GITR_SERVER, GITR_ORG_FORK, GITR_REPO_NAME),
		GITR_REPO_ABS_FSPATH:       fmt.Sprintf("%s/src/%s/%s/%s", os.Getenv("GOPATH"), GITR_SERVER, GITR_ORG_FORK, GITR_REPO_NAME), //$(GOPATH)/src/$(GITR_SERVER)/$(GITR_ORG_FORK)/$(GITR_REPO_NAME)
		GITR_BRANCH_NAME:           "master",
		repoController:             repo,
		GITR_COMMIT_MESSAGE:        " ",
	}, nil
}

func init() {
	var err error
	currentGitr, err = NewGitrAttrs()
	if err != nil {
		panic(err)
	}
}

// TODO:
// 	GITR_LAST_TAG:$(shell git describe --exact-match --tags $(shell git rev-parse HEAD))
// 	GITR_VERSION ?: $(shell echo $(TAGGED_VERSION) | cut -c 2-)
// 	GITR_COMMIT_MESSAGE ?: autocommit

// ---

// PrintAll prints gitr attributes
func PrintAll() error {

	gitrReflect := reflect.ValueOf(currentGitr)
	typeOfTarget := gitrReflect.Type()

	for i := 0; i < gitrReflect.NumField(); i++ {
		fmt.Printf("%s : %s \n", typeOfTarget.Field(i).Name, gitrReflect.Field(i).Interface())
	}

	return nil
}
