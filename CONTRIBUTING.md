# Contribution Guidelines

Please read this guide if you plan to contribute to this project. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues

If you find a bug while working with `gphotos-uploader-cli`, please open an issue on GitHub and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests

You are more than welcome to open issues in this project to suggest new features.

## Contributing Code

This project is mainly written in Golang.

> This project will maintain compatibility with the last two [golang major versions published](https://go.dev/doc/devel/release).

To contribute code:

1. Ensure you are running a supported golang version
1. Fork the project
1. Clone the project: `git clone https://github.com/[YOUR_USERNAME]/gphotos-uploader-cli && cd gphotos-uploader-cli`
1. Run `go mod download` to install the dependencies
1. Make changes to the code
1. Run `make build` to build the project
1. Make changes
1. Run tests and linter: `make ci`
1. Format your code: `go fmt ./...`
1. If users will notice the change, add an entry under `## Unreleased` in [CHANGELOG.md](CHANGELOG.md)
1. Commit changes
1. Push commits
1. Open pull request. The PR title must follow [Conventional Commits](https://www.conventionalcommits.org/), e.g. `fix: keep non-ASCII characters in album names`.

Maintainers: see [RELEASING.md](RELEASING.md) for how to publish a release.
