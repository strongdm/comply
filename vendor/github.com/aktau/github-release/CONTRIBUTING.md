## Releasing new versions

1) Bump the version in github/version.go

2) Add a commit with message "github-release v1.2.3"

3) Run `git tag v1.2.3` where "1.2.3" stands in for the version you actually
want.

4) Run `make release`. Be sure to have `GITHUB_TOKEN` set in your environment.
