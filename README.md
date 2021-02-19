![CI](https://github.com/amplify-edge/booty/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/alexadhy/booty/branch/master/graph/badge.svg?token=VLMYJWAQWJ)](https://codecov.io/gh/alexadhy/booty)

# (WIP) Booty

### Spec & Design: 

- [Issue 108](https://github.com/amplify-edge/main/issues/108)

## Usage (for devs)

1. Clone this repository
2. Run `make all`
3. Copy the binary in `bin` directory to your `PATH`
4. Run `booty`

## Tinker Tanker

To add a component, implement the `Component` 
interface in [here](https://github.com/alexadhy/booty/blob/master/dep/component.go)

This program installs config under:

- Linux: `$HOME/.local/booty/etc`
- Mac: `$HOME/Library/Application\ Support/booty/etc`
- Windows: `C:\\ProgramData\booty\etc`

The file `config.reference.json` is provided as a reference should you want to change the program behaviour.


## Implemented

1. Download 3rd party binaries and libraries our project needs
2. Install said binaries

## Short Term TODO

1. Also install service, managed by OS' service supervisor, preferably as user service on dev and as system-wide service running under unprivileged user on production.
2. Ability to run each service on the foreground like pingcap's `tiup playground`
3. Ability to self update and check for updates for third party binaries
4. Manage configs and backups for `maintemplatev2`
