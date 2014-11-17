
[![Build Status](https://drone.io/github.com/GeertJohan/ango/status.png)](https://drone.io/github.com/GeertJohan/ango/latest)
[![Issues](http://img.shields.io/github/issues/GeertJohan/ango.svg?style=flat-square)](https://github.com/GeertJohan/ango/issues)

## Ango: AngularJS <-> Go communication

`ango` is a tool that generates a protocol and API's for communication between [Go](http://golang.org) and [AngularJS](http://angularjs.org) over websockets.

**This project is under development and API's are likely to break.**

### Goals

The main goals are:
 - async server > client RPC.
 - async client > server RPC.
 - typed arguments and return values (see types below).
 - compatibility with Go packages `net/http` and `encoding/json`.
 - integrate into AngularJS as includable `module` providing a `service`.
 - underlying protocol and communication is not directly visible for user. Calls feel native/local.

What I don't want to do (not goals):
 - runtime discovery of available procedures.
 - A generic protocol designed for multiple languages/frameworks.

I chose to create a tool that generates Go and Angular/javascript without external dependencies, so the generated server and client code contain all information to communicate.

Generated Go code is a self-contained package without external imports. The generated code/package is to be imported by the application implementing/using the ango service.

For the client side a single `.js` file is generated containing an angular module. The module can be included by any other angular module.

### Terminology
A **service** exists of one or more **procedures** defined on the server- and/or client-side.
A **procedure** within a service is implemented on either the client- or server-side, and can be called by the other side.
A procedure can have zero or more arguments and zero or more return values.

### Notes
Several files describe the working of ango.

 - read [ango-definitions.md](notes/ango-definitions.md) about the .ango definition statements
 - read [types.md](notes/types.md) about the available types with ango procedures
 - read [protocol.md](notes/protocol.md) about the websockets/json protocol
 - read [thoughts.md](notes/thoughts.md) for idea's and upcomming features.

### Development

#### Devtool
This project involves several packages. There's a simple tool to automatically run and update the workspace as you go.
To run the tool, cd into the root folder (ango) and run: `sh tools/dev/run.sh`.
The devtool performs the following:
 - Watch ango cmd source (and imported packages) and re-build ango cmd on change.
 - Watch example source (and imported packages) and re-build on change.
 - Re-generate example ango service when ango tool was re-build or when .ango file changes or when a template changes.

#### Packages & directories
The ango source is divided into seperate packages:
 - `ango/definitions` contains types and strucutures defining an ango service. ([view godoc](http://godoc.org/github.com/GeertJohan/ango/definitions))
 - `ango/parser` implements a simple `.ango` definition file parser. The package provides functions and methods that are to be used directly by the generator and/or templates.
 - `ango` (main) is the cmd utilizing the above packages and contains the generators.

These packages exist to seperate logic and make it easier to create a more advanced parser (maybe using [yacc](http://golang.org/cmd/yacc/) and [nex](https://github.com/blynn/nex)).

Some other directories exist:
 - `ango/example` contains an example using ango. ([more info](example/README.md))
 - `ango/idea` old folder, marked for removal.
 - `ango/notes` documentation and ideas for this project.
 - `ango/templates` contains the templates used by the generators and are included by [go.rice](https://github.com/GeertJohan/go.rice).
 - `ango/tools/dev` contains the dev tool described above.
 - `ango/tools/publish` contains a tool ran by drone.io to build and preserve standalone binaries (linked in the download section below).

### TODO
There are still lots of things to do. Check out the [issues](https://github.com/GeertJohan/ango/issues) list. Please contact me if you would like to contribute.

### Download
At this time, automated builds are only available for linux_amd64. Download [release](https://drone.io/github.com/GeertJohan/ango/files/ango-release) for production. Or get the [latest](https://drone.io/github.com/GeertJohan/ango/files/ango-latest) build (nightly).

**Note on pretty javascript:** When [js-beautify](https://github.com/einars/js-beautify) is installed, it is used to clean up the generated javascript.

### License
This project is licensed under a Simplified BSD license. Please read the [LICENSE file](LICENSE).
