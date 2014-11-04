
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

### options such as direct channel
In some cases it would be nice to receive a channel from the returning function call go > angular. The return values and/or error would be sent on the channel. Maybe make this some kind of option keyword for procedure definitions.

### always return a result channel
For returning client procedures:

The ango internals in the Go server must use channels to communicate result data from the websocket to the blocking procedure call. It could be useful to expose the channel to the user. This allows the user to handle the result at a later time, or in certain situations skip the result all together (never read it from the channel).

It would be smart to make the channel buffered (cap 1) so the websocket handler won't block on sending on the channel.
