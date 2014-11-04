
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

