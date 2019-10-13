# Contribution Guidelines
Please read this guide if you plan to contribute to this project. We welcome any kind of contribution. No matter if you are an experienced programmer or just starting, we are looking forward to your contribution.

## Reporting Issues
If you find a bug while working with `gphotos-uploader-cli`, please [open an issue on GitHub](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?assignees=pacoorozco&labels=bug&template=bug_report.md) and let us know what went wrong. We will try to fix it as quickly as we can.

## Feature Requests
You are more than welcome to open issues in this project to [suggest new features](https://github.com/gphotosuploader/gphotos-uploader-cli/issues/new?assignees=&labels=feature+request&template=feature_request.md).

## Contributing Code
This project is mainly written in Golang.

> This project will maintain compatibility with the last two Go major versions published. Currently Go 1.12 and Go 1.13.

To contribute code:
1. Ensure you are running golang version 1.12 or greater
2. Set the following environment variables:
    ```
    GO111MODULE=on
    ```
3. Fork the project
4. Clone the project: `git clone https://github.com/[YOUR_USERNAME]/gphotos-uploader-cli && cd gphotos-uploader-cli`
5. Run `go mod download` to install the dependencies
6. Make changes to the code
7. Run `make build` to build the project
8. Make changes
9. Run tests: `make test`
10. Format your code: `go fmt ./...`
11. Commit changes
12. Push commits
13. Open pull request
