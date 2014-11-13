
### JS procedure call type parsing
**Currently:** Arguments used in a procedure call must match the js type (string, number) for the procedure arg params. This requires the user to parse input fields. If an argument value does not match the required type, an error is thrown.

**Idea:** Let the generated client code try to parse procedure arguments before asserting their type. For instance, `add('1', 2)` will parse `'1'` to `1` using parseInt(). And `notify(42)` will parse `42` to `'42'` using parseString. And `sendString(object)` will try to convert `object` to string using `object.toString()`.

### Faster JSON
Generated structures marshal a lot faster with ffjson: https://github.com/pquerna/ffjson

### service blocks and includable files
Includable files would be a benefit for large project where custom types are used by different services. When writing this, a single ango file defines a single service with the name specified at the top. This is not an ideal setup when including an ango file into another file. Syntax could be changed to something like:
```
type customStringType string

service myService {
	server add(a int32, b int32) (c int32)
	client notify(text customStringType)
}
```

An include statement would then add all types from an other file:

common.ango:
```
type customString string
type customInt32 int32
type foobar {
	foo string
	bar int32
}
```

myService.ango:
```
include "common.ango"

service {
	server oneway addFoobar(fb foobar)
}
```

### ~~options such as direct channel~~
~~In some cases it would be nice to receive a channel from the returning function call go > angular. The return values and/or error would be sent on the channel. Maybe make this some kind of option keyword for procedure definitions.~~

### always return a result channel
For returning client procedures:

The ango internals in the Go server must use channels to communicate result data from the websocket to the blocking procedure call. It could be useful to expose the channel to the user. This allows the user to handle the result at a later time, or in certain situations skip the result all together (never read it from the channel).

It would be smart to make the channel buffered (cap 1) so the websocket handler won't block on sending on the channel.

### generated Go code as a package
This describes a new approach to project setup and ango cmd usage.

Standardize `name foobarservice` (all lowercase).

In project have:
```
/myproject
	/chatservice          // folder containing generated code for chatservice.ango, package chatservice
		/types.gen.go     // generated go code for custom types defined in chatservice.ango
		/server.gen.go    // generated go code implementing the server as defined in chatservice.ango
	/chatservice.ango     // ango definitions for chatservice
	/chatservice.gen.js   // generated js code implementing types and client as defined in chatservice.ango
	/main.go              // main program code
```

main.go could contain `// go:generate` clause.

Running ango would be as simple as running `ango chatservice.ango` in the working directory `myproject`.

This would generate go sources in the myproject/chatservice folder. Javascript sources are generated to /myproject/chatservice.gen.js

Optionally a different path for the javascript or go can be given with the --js-path and --go-path flags which are either relative to ango's working directory or absolute.

The package myproject/chatservice can contain custom (non-generated) code with access to the service internals. This could be used to create hooks (next topic)

### Hooks
It could be useful to have access to the service internal (debugging, low level access).

This could be achieved by doing:

```
// generated code:
var someHook func(val string)

func serviceInternals(val string) {
	if someHook != nil {
		someHook(val)
	}
}
```

Custom code could implement this:
```
func init() {
	someHook = func(val string) {
		fmt.Printf("interesting: %s\n", val)
	}
}
```

### Javascript procedure: return promise

When a javascript procedure implementation returns a promise it will defer the callback to the promise resolve.

### Javascript procedure: deferred implementation

A javascript procedure handler can't always be set during the AngularJS config phase. There should be an option to defer setting the handler until a later moment. Incomming calls will be buffered and played back when the handler is set. Maybe do something like:
```
chatserviceProvider.setHandlers({
	askQuestion: 'defer', 
});
```

There could be multiple options, e.g.: 'defer' and 'defer-sync'. 'defer' will let ango make all buffered calls directly when the handler is set (all at once), whereas 'defer-sync' will play the calls one by one (in received order). So when 2 calls were buffered, it will only make the second call when the first returned. See idea below about buffered calls..

### Buffering of calls or: sequential procedures

Above is described how several buffered calls to a handler that was not yet configured will be called all at once when it is set, or one by one (sequential). Maybe this behaviour should be configured in the ango definition file, and it could work for both server and client: buffered procedures. A sequential procedure will handle only one call at once. 'sequential' could be a keyword. e.g.:
```
server sequential askQuestion(question string) (answer string)
server sequential oneway notify(message string)
```

