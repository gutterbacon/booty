package gitr

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

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
func ForkSetupOld() error {

	// 	## Sets up the git fork locally.
	// gitr-fork-setup-old:
	// 	# Pre: you git forked ( via web) and git cloned (via ssh)
	// 	# add upstream repo

	// 	#git remote add upstream git://$(GITR_SERVER)/$(GITR_ORG_UPSTREAM)/$(GITR_REPO_NAME).git

	return errors.New("Not implemented")

}

// ForkSetup sets up the git fork locally.
func ForkSetup() error {

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

	fmt.Println("Pre: you git forked ( via web) and git cloned (via ssh)")
	fmt.Println("Sets up git config upstreak to point to the upstream origin")
	fmt.Println()
	fmt.Println("EX git remote add upstream git@github.com-joe-getcouragenow:getcouragenow/dev")
	fmt.Printf("EX git remote add upstream git@%s-%s:%s/%s \n", currentGitr.GITR_SERVER, currentGitr.GITR_USER, currentGitr.GITR_ORG_UPSTREAM, currentGitr.GITR_REPO_NAME)

	_, err := currentGitr.repoController.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{fmt.Sprintf("git@%s-%s:%s/%s", currentGitr.GITR_SERVER, currentGitr.GITR_USER, currentGitr.GITR_ORG_UPSTREAM, currentGitr.GITR_REPO_NAME)},
	})

	return err

}

// ForkCatchup Sync upstream with your fork. Use this to make a PR.
func ForkCatchup() error {

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

	// err := currentGitr.repoController.Fetch(&git.FetchOptions{
	// 	RemoteName: "upstream",
	// })

	wt, err := currentGitr.repoController.Worktree()
	if err != nil {
		return err
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName: "upstream",
	})

	return err

}

// ForkCommit Commit the changes to the repo
func ForkCommit() error {

	// 	## Commit the changes to the repo
	// gitr-fork-commit:
	// 	@echo GITR_COMMIT_MESSAGE: $(GITR_COMMIT_MESSAGE)
	// 	git add --all
	// 	git commit -m '$(GITR_COMMIT_MESSAGE)'

	wt, err := currentGitr.repoController.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Add(".")

	if err != nil {
		return err
	}

	wt.Commit(currentGitr.GITR_COMMIT_MESSAGE, &git.CommitOptions{
		All: true,
	})

	if err != nil {
		return err
	}

	return nil

}

// ForkPush Push the repo to orgin
func ForkPush() error {

	// 	## Push the repo to orgin
	// gitr-fork-push:
	// 	git push origin $(GITR_BRANCH_NAME)

	return currentGitr.repoController.Push(&git.PushOptions{
		RemoteName: "origin",
	})

}

// ForkOpen Opens the forked git server.
func ForkOpen() error {

	// 	## Opens the forked git server.
	// gitr-fork-open:
	// 	open $(GITR_REPO_ABS_URL).git

	// TODO : I don't get what it's really supposed to do
	return errors.New("Not implemented")
}

// ForkPRSubmit Submits the PR you pushed
func ForkPRSubmit() error {
	// ## Submits the PR you pushed
	// gitr-fork-pr-submit:
	// 	## TODO. Alex gave me the commands.
	// 	open $(GITR_REPO_ABS_URL).git

	return errors.New("Not implemented")

}
