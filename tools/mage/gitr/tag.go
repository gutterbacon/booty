package gitr

import "errors"

// TagCreate creates a tag.
func TagCreate() error {

	// ## Create a tag.
	// gitr-tag-create:
	// 	# this will create a local tag on your current branch and push it to Github.

	// 	git tag $(GIT_TAG_NAME)

	// 	# push it up
	// 	git push origin --tags

	// TODO: Need to setup GIT_TAG_NAME first...

	return errors.New("Not implemented")

}

// TagDelete is deletes a tag
func TagDelete() error {

	// 	## Deletes a tag.
	// gitr-tag-delete:
	// 	# this will delete a local tag and push that to Github

	// 	git push --delete origin $(GIT_TAG_NAME)
	// 	git tag -d $(GIT_TAG_NAME)

	// TODO: Need to setup GIT_TAG_NAME first...

	return errors.New("Not implemented")

}
