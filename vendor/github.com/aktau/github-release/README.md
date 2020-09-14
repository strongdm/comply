github-release
==============

A small commandline app written in Go that allows you to easily create
and delete releases of your projects on Github. In addition it allows
you to attach files to those releases.

It interacts with the [github releases API][releases-api]. Though it's entirely
possible to [do all these things with cURL][curl], it's not really that
user-friendly. For example, you need to first query the API to find the id of
the release you want, before you can upload an artifact. `github-release` takes
care of those little details.

[curl]: https://github.com/blog/1645-releases-api-preview
[releases-api]: https://developer.github.com/v3/repos/releases

It might still be a bit rough around the edges, pull requests are
welcome!

How to install
==============

If you don't have the Go toolset installed, and you don't want to, but
still want to use the app, you can download binaries for your platform
on the [releases
page](https://github.com/github-release/github-release/releases/latest). Yes, that's
dogfooding, check the makefile!

If you have Go installed, you can just do:

```sh
go get github.com/github-release/github-release
```

This will automatically download, compile and install the app.

After that you should have a `github-release` executable in your
`$GOPATH/bin`.

How to use
==========

**NOTE**: for these examples I've [created a github token][token] and set it as
the env variable `GITHUB_TOKEN`. `github-release` will automatically pick it up
from the environment so that you don't have to pass it as an argument.

[token]: https://help.github.com/articles/creating-an-access-token-for-command-line-use

```sh
# set your token
export GITHUB_TOKEN=...

# check the help
$ github-release --help

# make your tag and upload
$ git tag ... && git push --tags

# check the current tags and existing releases of the repo
$ github-release info -u aktau -r gofinance
git tags:
- v0.1.0 (commit: https://api.github.com/repos/aktau/gofinance/commits/f562727ce83ce8971a8569a1879219e41d56a756)
releases:
- v0.1.0, name: 'hoary ungar', description: 'something something dark side 2', id: 166740, tagged: 29/01/2014 at 14:27, published: 30/01/2014 at 16:20, draft: ✔, prerelease: ✗
  - artifact: github.go, downloads: 0, state: uploaded, type: application/octet-stream, size: 1.9KB, id: 68616

# create a formal release
$ github-release release \
    --user aktau \
    --repo gofinance \
    --tag v0.1.0 \
    --name "the wolf of source street" \
    --description "Not a movie, contrary to popular opinion. Still, my first release!" \
    --pre-release

# you've made a mistake, but you can edit the release without
# having to delete it first (this also means you can edit without having
# to upload your files again)
$ github-release edit \
    --user aktau \
    --repo gofinance \
    --tag v0.1.0 \
    --name "Highlander II: The Quickening" \
    --description "This is the actual description!"

# upload a file, for example the OSX/AMD64 binary of my gofinance app
$ github-release upload \
    --user aktau \
    --repo gofinance \
    --tag v0.1.0 \
    --name "gofinance-osx-amd64" \
    --file bin/darwin/amd64/gofinance

# upload other files...
$ github-release upload ...

# you're not happy with it, so delete it
$ github-release delete \
    --user aktau \
    --repo gofinance \
    --tag v0.1.0
```

GitHub Enterprise Support
=========================
You can point to a different GitHub API endpoint via the environment variable ```GITHUB_API```:

```
export GITHUB_API=http://github.company.com/api/v3
```

Used libraries
==============

| Package                                                                    | Description         | License |
| ------------------------------------------------------------------------   | ------------------- | ------- |
| [github.com/dustin/go-humanize](https://github.com/dustin/go-humanize)     | humanize file sizes | MIT     |
| [github.com/tomnomnom/linkheader](https://github.com/tomnomnom/linkheader) | GH API pagination   | MIT     |
| [github.com/voxelbrain/goptions](https://github.com/voxelbrain/goptions)   | option parsing      | BSD     |
| [github.com/kevinburke/rest](https://github.com/kevinburke/rest)   | HTTP client      | MIT     |

Todo
====

- Check if an artifact is already uploaded before starting a new upload

Copyright
=========

Copyright (c) 2014-2017, Nicolas Hillegeer. All rights reserved.
Copyright (c) 2020, Meter, Inc. All rights reserved.
