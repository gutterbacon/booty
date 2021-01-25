# Booty

Inspired by https://github.com/pingcap/tidb-operator and https://github.com/pingcap/tiup

Booty installs the dependencies for your local machine, for the CI ( github actions ), or for your deployment ot a laptop or sever with zero changes to the make files or a standard CI GitHub script.

It means that the CI github actions call the same makefile that a developer calls on their local machine. CI calls makefile that then calls booty. Local dev just calls make file directly. It’s all the same.

The only thing the CI git hbu action script installs is golang and flutter. So it’s same as dev local where it’s expected you have golang, flutter, make and git installed but nothing else because booty does all the other dependency installs for you.

For Users, it also can deploy any things needed on a server outside our main cli and server.

The real code that does all the work is in shared, and so booty imports shared. Its designed this way so that main Cli and Server also import shared as needed.

So from booty you can install everything you need.

Booty installs any binaries needed by a dev or user and the tools used by a dev or user.

Each repo has the same CI script, that just uses a github actions to install go and flutter and then calls the makefile....


https://github.com/getcouragenow/booty/blob/master/.github/workflows/ci.yml is enough for all repos to do what they need to do.

The CI installs golang and then calls the make file target called "all".