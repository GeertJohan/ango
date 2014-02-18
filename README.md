## Ango: Angular <-> Go communication

`ango` is a tool that generates a protocol for communication between [Go](http://golang.org) and [AngularJS](http://angularjs.org) over http/websockets.

**This project is still under development, please do look arround and let me know if you have good idea's**

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
This project involves several packages. There's a simple tool to automatically run and update the workspace as you go.
This tool depends on `rerun` and `go.sgr`. Install them using `go get github.com/skelterjohn/rerun` and `go get github.com/foize/go.sgr`.
To run the tool, cd into the root folder (ango) and run: `go run tools/dev.go`.
This performs the following:
 - Watch ango source (and imported packages) and re-build on change.
 - Watch example source (and imported packages) and re-build on change.
 - Re-genereate example ango service when ango tool was re-build or when .ango file changes.

### License
This project is licensed under a Simplified BSD license. Please read the [LICENSE file](LICENSE).