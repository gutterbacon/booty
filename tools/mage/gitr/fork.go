package gitr

import "fmt"

// ForkAll executes all commands down there
func ForkAll() {
	// gitr-fork-all: gitr-status gitr-fork-catchup gitr-status gitr-fork-commit gitr-fork-push gitr-fork-open
	// #See: https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/syncing-a-fork
}

// ForkCloneTemplate prints the exact Git Clone command you need :)
func ForkCloneTemplate() error {

	// 	# You can call this from your Org folder:
	// 	# make -f shared/boilerplate/gitr.mk gitr-fork-clone-template
	fmt.Println(`
Template is:
EX "git clone git@github.com-ME-getcouragenow:ME-getcouragenow/REPO_NAME"

So if your fork is:
github.com/james-getcouragenow/dev

You use:
EX: "git clone git@github.com-james-getcouragenow:james-getcouragenow/dev"
	`)

	return nil

}

// ForkSetupOld : [DEPRECATED] : Sets up the git fork locally.
func ForkSetupOld() {

	// 	## Sets up the git fork locally.
	// gitr-fork-setup-old:
	// 	# Pre: you git forked ( via web) and git cloned (via ssh)
	// 	# add upstream repo

	// 	#git remote add upstream git://$(GITR_SERVER)/$(GITR_ORG_UPSTREAM)/$(GITR_REPO_NAME).git

}

// ForkSetup sets up the git fork locally.
func ForkSetup() {

	// 	## Sets up the git fork locally.
	// gitr-fork-setup:
	// 	# Pre: you git forked ( via web) and git cloned (via ssh)
	// 	# Sets up git config upstreak to point to the upstream origin
	// 	@echo
	// 	@echo EX git remote add upstream git@github.com-joe-getcouragenow:getcouragenow/dev
	// 	@echo
	// 	@echo EX git remote add upstream git@$(GITR_SERVER)-$(GITR_USER):$(GITR_ORG_UPSTREAM)/$(GITR_REPO_NAME)
	// 	@echo
	// 	# WORKS
	// 	git remote add upstream git@$(GITR_SERVER)-$(GITR_USER):$(GITR_ORG_UPSTREAM)/$(GITR_REPO_NAME)
	// 	@echo

}

// ForkCatchup Sync upstream with your fork. Use this to make a PR.
func ForkCatchup() {

	// 	## Sync upstream with your fork. Use this to make a PR.
	// gitr-fork-catchup:

	// 	# This fetches the branches and their respective commits from the upstream repository.
	// 	@echo
	// 	git fetch upstream
	// 	@echo

	// 	# This brings your fork's master branch into sync with the upstream repository, without losing your local changes.
	// 	@echo
	// 	git merge upstream/$(GITR_BRANCH_NAME)
	// 	@echo

}

// ForkCommit Commit the changes to the repo
func ForkCommit() {

	// 	## Commit the changes to the repo
	// gitr-fork-commit:
	// 	@echo GITR_COMMIT_MESSAGE: $(GITR_COMMIT_MESSAGE)
	// 	git add --all
	// 	git commit -m '$(GITR_COMMIT_MESSAGE)'

}

// ForkPush Push the repo to orgin
func ForkPush() {

	// 	## Push the repo to orgin
	// gitr-fork-push:
	// 	git push origin $(GITR_BRANCH_NAME)

}

// ForkOpen Opens the forked git server.
func ForkOpen() {

	// 	## Opens the forked git server.
	// gitr-fork-open:
	// 	open $(GITR_REPO_ABS_URL).git

}

// ForkPRSubmit Submits the PR you pushed
func ForkPRSubmit() {
	// ## Submits the PR you pushed
	// gitr-fork-pr-submit:
	// 	## TODO. Alex gave me the commands.
	// 	open $(GITR_REPO_ABS_URL).git
}
