## Ango: Angular <-> Go communication

`ango` is a tool that generates a protocol for communication between [Go](http://golang.org) and [AngularJS](http://angularjs.org) over http/websockets.

**This project is still under development, please do look around and let me know if you have good ideas**

### Goals

The main goals are:
 - async server > client RPC.
 - async client > server RPC.
 - typed arguments and return values (see types below).
 - work with Go packages `net/http` and `encoding/json`.
 - integrate into AngularJS as includable module + ng service.
 - underlying protocol and communication is not directly visible for user. Calls feel native/local.

What I don't want to do:
 - runtime discovery of available procedures.
 - A generic protocol designed for multiple languages/frameworks.

Therefore I chose to create a tool that generates Go and Angular/javascript, so both server and client hold all information to communicate and know whats comming at any time.

Code generated for Go can be copied into any go package. The code doesn't form a package itself.
NEEDS THINKING: this has drawbacks on the generated code (need to use scope to hide inner variables from the rest of package).

For angular a single `.js` file is generated  holding an angular module. The module can be included by any other angular module.

### Terminology
A **service** exists of one or more **procedures** server- and/or client-side.
A **procedure** within a service resides client- or server-side, and can be called by the other party.
A procedure can have zero or more arguments and zero or more return values.

### Notes
Several files describe the working of ango.

 - read [ango-definitions.md](notes/ango-definitions.md) about the .ango definition statements
 - read [types.md](notes/types.md) about the available types with ango procedures
 - read [protocol.md](notes/protocol.md) about the websockets/json protocol

### Development

[![Build Status](https://drone.io/github.com/GeertJohan/ango/status.png)](https://drone.io/github.com/GeertJohan/ango/latest)

#### Devtool
This project involves several packages. There's a simple tool to automatically run and update the workspace as you go.
To run the tool, cd into the root folder (ango) and run: `sh tools/dev/run.sh`.
The devtool performs the following:
 - Watch ango source (and imported packages) and re-build on change.
 - Watch example source (and imported packages) and re-build on change.
 - Re-generate example ango service when ango tool was re-build or when .ango file changes.

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
There's still lots of things to do. If you wish to help out, please contact me.

### Download
At this time, automated builds are only available for linux_amd64. Download [release](https://drone.io/github.com/GeertJohan/ango/files/ango-release) for production. Or get the [latest](https://drone.io/github.com/GeertJohan/ango/files/ango-latest) build (nightly).

**Note on pretty javascript:** When [js-beautify](https://github.com/einars/js-beautify) is installed, it is used to clean up the generated javascript.

### License
This project is licensed under a Simplified BSD license. Please read the [LICENSE file](LICENSE).
