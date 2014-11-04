## Ango definitions
The ango tool takes a file with ango definitions to generate the communication protocol. This file describes the syntax for these definitions.


The `.ango` definition file  specifies the protocol name, and several services.

### Comments
Comments can be placed on any line and are started with `//`. Everything until newline (`\n`) is ignored.

### Letters and digits

```
letter = "a" … "z" | "A" … "Z" .
digit = "0" … "9" .
```

### Identifier
Identifiers in ango syntax are not entirely equal to Go identifiers. Ango identifiers only allow `a-z` and `A-Z` letters.

`identifier = letter { letter | digit } .`

Some identifiers are predeclared.

### Service name
The first actual statement should be the `name` statement:

```
ServiceClause = "name" ServiceName .
ServiceName = identifier .
```

It is good practice to match the filename. e.g. `name myService` for `myService.ango`.

### Types
Type description:

```
Type      = TypeName | TypeLit .
TypeName  = identifier .
TypeLit   = SliceType | MapType | StructType .
```

#### Slice types
A slice type translates to a `slice` in Go. In javascript it is translated to an array. Defining the length or capacity of elements is not possible at this moment.

```
SliceType    = "[" "]" ElementType .
ElementType  = Type .
```

#### Map types
A map translates to a `map` in Go. In javascript this is translated to an object.

```
MapType  = "map" "[" KeyType "]" ElementType .
KeyType  = Type .
```

#### Struct types
A struct type is directly compatible Go code. In javascript this is represented as a 

```
StructType  = "struct" "{" { FieldDecl } "}" .
FieldDecl   = identifier Type .
```

#### Builtin types
Ango provides a set of builtin types such as integers and strings. Because javascript has less types than Go, different Go types translate to the same javascript type. For instance Go's `uint8`, `uint64`, `int8` and `int32` all translate to javascripts `number`. Read more details about the builtin types in [types.md](types.md).

### Type declarations:
Custom types can be declared within a ango file.

```
TypeDecl  = "type" TypeSpec .
TypeSpec  = identifier Type .
```

### Procedure
Procedure description:

```
ProcedureDecl                = ServerProcedureSpec | ClientProcedureSpec .
ServerProcedureSpec          = "server" (OnewayProcedureSpec|ReturningProcedureSpec) .
ClientProcedureSpec          = "client" (OnewayProcedureDecl|ReturningProcedureSpec) .
OnewayProcedureSignature     = "oneway" ProcedureName Parameters .
ReturningProcedureSignature  = ProcedureName Parameters [ Result ] .
ProcedureName                = identifier .
Result                       = Parameters .
Parameters                   = "(" [ ParameterList ] ")" .
ParameterList                = ParameterDecl { "," ParameterDecl } .
ParameterDecl                = identifier Type .
```

The first keyword, `'server'` or `'client'`, indicates which side provides/implements the procedure.

A `oneway` procedure does not wait for a response from the other side. A call to a `oneway` procedure returns immediately after the call has been sent over the websocket. There's no result possible. Any possible error should be handled at the called side only, as none can be sent back. 

A `returning` procedure call retuns when the procedure implementation has returned (with or without error). Optionally, some return values can be sent back.

Identifiers (`serviceName`, `procedureName`, `argName`, `retName`) must all be in lowerCamelCase.

### Example
There's an example `.ango` file at [/example/example.ango](/example/example.ango)


### Idea's
 - appendix `notifies(typeA)` for a server procedure, which adds a notify argument to the server procedure handler. The notify argument (`type func(typeA)`) can be called by the procedure implementation.