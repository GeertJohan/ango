## Ango definitions
The ango tool takes a file with ango definitions to generate the communication protocol. This file describes the syntax for these definitions.


The `.ango` definition file specifies the protocol name, and several services.

Comments can be placed on any line and are started with `//`. Everything until newline (`\n`) is ignored.

`// this is a comment`

Name, first non-comment/non-empty line:

`name <serviceName>`

Procedure description:

`{'server'|'client'} ['oneway'] procedureName '(' argument [',' argument]* ')' [ '('result [',' result]* ')' ]`

 - `server`/`client` indicates which party provides the procedure.
 - `synchronized` (idea) see Idea's section.
 - `oneway` (idea) indicates a fire-and-forget procedure. The caller returns imediatly once the call has been sent over the websocket. There's no result expected back. Any possible error should be handled server-side only. Cannot be combined with the `synchronized` keyword.
 - `args` is a list of argument names and their type.
 - `rets` is a list of return values and their type. `rets` is not available for oneway procedures.

Argument/result:
`name type`

