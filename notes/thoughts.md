
### JS procedure call type parsing
**Currently:** Arguments used in a procedure call must match the js type (string, number) for the procedure arg params. This requires the user to parse input fields. If an argument value does not match the required type, an error is thrown.

**Idea:** Let the generated client code try to parse procedure arguments before asserting their type. For instance, `add('1', 2)` will parse `'1'` to `1` using parseInt(). And `notify(42)` will parse `42` to `'42'` using parseString. And `sendString(object)` will try to convert `object` to string using `object.toString()`.