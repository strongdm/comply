![Comply](https://github.com/strongdm/comply/blob/master/logo.png)

Comply is a SOC2-focused compliance automation tool:

- **Policy Generator**: markdown-powered **document pipeline** for publishing auditor-friendly **policy documents**
- **Ticketing Integration**: automate compliance throughout the year via your existing **ticketing system**
- **SOC2 Templates**: open source policy and procedure **templates** suitable for satisfying a SOC2 audit

# Installation

macOS:

`brew tap strongdm/comply; brew install comply`

Linux:

[Download latest release](https://github.com/strongdm/comply/releases)

Go users:

`go get github.com/strongdm/comply`

# Get Started

Start with `comply init`:

```
$ mkdir my-company
$ cd my-company
$ comply init
```

Once `comply init` is complete, just `git init` and `git push` your project to a new repository. You're ready to begin editing the included policy boilerplate text.

# Discussion

Join us in [Comply Users](https://join.slack.com/t/comply-users/shared_invite/enQtMzU3MTk5MDkxNDU4LTMwYzZkMjA4YjQ2YTM5Zjc0NTAyYWY5MDBlOGMwMzRmZTk5YzBlOTRiMTVlNGRlZjY1MTY1NDE0MjY5ZjYwNWU)

# Screenshots

[Demo video](https://vimeo.com/270257486)

## Start a Project
![screencast 1](sc-1.gif)

## Build PDFs
![screencast 4](sc-4.gif)
![pdf example](pdf-example.png)


## Track Policy Coverage
![screencast 3](sc-2.gif)

## Dashboard
![screencast 2](sc-3.gif)

## Dependencies

Comply relies on [pandoc](https://pandoc.org/), which can be installed directly as an OS package or invoked via Docker.

## CLI

```
NAME:
   comply - policy compliance toolkit

USAGE:
   comply [global options] command [command options] [arguments...]

COMMANDS:
     init             initialize a new compliance repository (interactive)
     build, b         generate a static website summarizing the compliance program
     procedure, proc  create ticket by procedure ID
     scheduler        create tickets based on procedure schedule
     serve            live updating version of the build command
     sync             sync ticket status to local cache
     todo             list declared vs satisfied compliance controls
     help, h          Shows a list of commands or help for one command
```
## Ticketing Integrations:
- Jira
- Github
- Gitlab

## Forking and local development
> Assumes installation of golang and configuration of GOPATH in .bash_profile, .zshrc, etc
> Inspiration: http://code.openark.org/blog/development/forking-golang-repositories-on-github-and-managing-the-import-path

```
$ go get http://github.com/strongdm/comply
$ cd $GOPATH/src/github.com/strongdm/comply ; go get ./...
$ make
$ cd example
$ mv comply.yml.example comply.yml
$ ../comply -h
$ ../comply sync
$ ../comply serve
#
$ make # recompile as needed with in $GOPATH/src/github.com/strongdm/comply
```
