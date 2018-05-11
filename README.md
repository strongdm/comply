# Comply

Comply is a SOC2-focused compliance automation tool. Comply features a markdown-powered **document pipeline** and a git-powered **workflow** that help policies and procedures _feel_ like software development.

Comply manages the lifecycle of your program throughout the year via your existing **ticketing system**.

In addition to automation, Comply includes a SOC2-focused module featuring open source policy and procedure **templates** suitable for satisfying a SOC2 audit.

# Getting Started

macOS:

`brew install comply`

Go users:

`go get github.com/strongdm/comply`

# Discussion

Join us in [Comply Users](https://join.slack.com/t/comply-users/shared_invite/enQtMzU3MTk5MDkxNDU4LTMwYzZkMjA4YjQ2YTM5Zjc0NTAyYWY5MDBlOGMwMzRmZTk5YzBlOTRiMTVlNGRlZjY1MTY1NDE0MjY5ZjYwNWU)

# Screenshots

## Start a Project
![screencast 1](sc-1.gif)

## Track Policy Coverage
![screencast 3](sc-2.gif)

## Dashboard
![screencast 2](sc-3.gif)

## CLI

```
NAME:
   comply - policy compliance toolkit

USAGE:
   comply [global options] command [command options] [arguments...]

COMMANDS:
     build, b   generate a static website summarizing the compliance program
     init       initialize a new compliance repository (interactive)
     scheduler  create tickets based on procedure schedule
     serve      live updating version of the build command
     sync       sync external systems to local data cache
     todo       list declared vs satisfied compliance controls
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```