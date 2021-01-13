package gitr

import (
	"errors"
)

// ReleaseTag stages a release (usage: make release-tag VERSION={VERSION_TAG})
func ReleaseTag() error {

	// ## Stage a release (usage: make release-tag VERSION={VERSION_TAG})
	// gitr-release-tag:
	// 	@echo Tagging release with version "${VERSION}"
	// 	@git tag -a ${VERSION} -m "chore: release version '${VERSION}'"
	// 	@echo Generating changelog
	// 	@git-chglog -o CHANGELOG.md
	// 	@git add CHANGELOG.md
	// 	@git commit -m "chore: update changelog for version '${VERSION}'"

	// TODO : need to setup VERSION first on gitrAttributes and New..

	return errors.New("Not implemented")

}

// ReleasePush pushes a release (warning: ensure the release was staged first)
func ReleasePush() error {
	// 	## Push a release (warning: ensure the release was staged first)
	// gitr-release-push:
	// 	@echo Publishing release
	// 	@git push --follow-tags

	return errors.New("Not implemented")

}
