# Booty

Inspired by https://github.com/pingcap/tidb-operator

Booty installs the dependencies for your local machine or for the CI ( github actions ) with zero changes to the make files or a standard CI GitHub script.

It maeans that the CI github actions call the same makefile that a developer calls on their local machine. CI calls makefile that then calls booty. Local dev just calls make file directly. It’s all the same.

The only thing CI script installs is golang anf flutter. So it’s same as dev local where it’s expected you have golang, flutter, make and git installed but nothing else because booty does all the other dependency installs for you.

The real code in in shared, and so booty import shared. Its designed this way so that Cli and Server also import shared as needed.

So from booty you can install everything you need.

Booty installs any binaries needed by a dev or user and the tools used by a dev or user.

Each repo has the same CI script, that just uses a github actions to install go and flutter and then calls the makefile....