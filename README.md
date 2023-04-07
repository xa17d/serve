# serve

A simple local webserver that serves files of a directory.

By default, a small piece of JavaScript is injected into HTML pages to automatically reload when any file within the directory changes.

_Only tested on MacOS, but should run on Linux and Windows too._

## Usage


```sh
cd to/the/path/you/want/to/serve
serve
```

### Options

| Option         | Default          | Description                                                                                                                                                     |
|----------------|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `address`      | `localhost:3000` | Address of the server.<br><br>**Example:**<br>`serve --address="localhost:8080"`                                                                                |
| `folder`       | `.`              | Folder that should be served.<br><br>**Example:**<br>`serve --folder="/my/folder"`                                                                              |
| `auto-refresh` | `true`           | Inject JavaScript into HTML pages that automatically reloads when a file in the specified folder changes.<br><br>**Example:**<br>`serve --auto-refresh="false"` |
| `open`         | `true`           | Opens the server address in the default browser on startup.<br><br>**Example:**<br>`serve --open="false"`                                                       |

## Installation

No installation required, but recommended to copy `serve` to a folder that is added to `$PATH`, e.g. `/usr/local/bin`.

You can find the binaries on the [Releases page](https://github.com/xa17d/serve/releases)

## Tech Stack

A single [Go](https://go.dev/) file.

Based on the Go standard library's [net/http](https://pkg.go.dev/net/http).

## Build

**Prerequisites:** Go 1.19 _(or higher)_

To build and debug locally, run `go run .`.

To build for all platforms, run `./build.sh`.
Result will be put into `build` folder.

## Release

_(This is only relevant for maintainers)_

1. Increment `version` variable in [main.go](main.go).
2. Run `./build.sh`
3. Commit & Push
4. Create new release on GitHub and attach zip files from the `build` folder.

## License

Distributed under [BSD 3-Clause License](LICENSE). See [LICENSE](LICENSE) for more information.

## Contact

- Website: [xa1.at/serve](https://xa1.at/serve/)
- E-Mail: [support@xa1.at](mailto:support@xa1.at?subject=serve%20app)