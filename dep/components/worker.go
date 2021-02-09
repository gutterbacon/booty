package components

import (
	"go.amplifyedge.org/booty-v2/dep"
)

type workerType struct {
	jobType      string
	targetDir    string
	dependencies []dep.Component
	errChan      chan error
}

func newWorkerType(jobType string, targetDir string, deps []dep.Component, errChan chan error) *workerType {
	return &workerType{
		jobType:      jobType,
		targetDir:    targetDir,
		dependencies: deps,
		errChan:      errChan,
	}
}

func (w *workerType) do(i int) {
	switch w.jobType {
	case "download":
		w.errChan <- w.dependencies[i].Download(w.targetDir)
	case "install":
		w.errChan <- w.dependencies[i].Install()
	case "uninstall":
		w.errChan <- w.dependencies[i].Uninstall()
	}
}
