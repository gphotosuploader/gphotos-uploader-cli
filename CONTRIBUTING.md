# Contribution Guidelines
Please read this guide if you plan to contribute to this project. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with `gphotos-uploader-cli`, please [open an issue on GitHub](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?assignees=pacoorozco&labels=bug&template=bug_report.md) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?assignees=&labels=feature+request&template=feature_request.md).

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
1. Run tests: `make test`
1. Run linter: `make lint`
1. Format your code: `go fmt ./...`
1. Commit changes
1. Push commits
1. Open pull request
