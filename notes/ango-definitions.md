## Ango definitions
The ango tool takes a file with ango definitions to generate the communication protocol. This file describes the syntax for these definitions.


The `.ango` definition file  specifies the protocol name, and several services.

### Comments
Comments can be placed on any line and are started with `//`. Everything until newline (`\n`) is ignored.

### Service name
The first actual statement should be the `name` statement:
`name <serviceName>`

### Procedure
Procedure description:
`('server'|'client') ['oneway'] procedureName '(' [ argName argType {',' argName argType} ] ')' [ '(' retName retType { ',' retName retType } ')' ]`

The first keyword, `'server'` or `'client'`, indicates which party provides the procedure.
When the keyword `'oneway'` is added, the caller returns imediatly once the call has been sent over the websocket. There's no result expected. Any possible error should be handled at the called side only, as none can be sent back. The `procedureName` is used throughout the generated code and is exposed to the users custom code. Following the `procedureName` are one or two groups enclosed by parenthesis. The first for arguments. The second for return values. The arguments group can be empty. The return values group can not, and should be omitted al together when no values are to be returned by the procedure. A oneway procedure can not have return values.

Identifiers (`serviceName`, `procedureName`, `argName`, `retName`) must all be in lowerCamelCase.

Available types (`argType` and `retType`) are defined in [types.md](types.md).

### Example
There's an example `.ango` file at [/example/example.ango](/example/example.ango)