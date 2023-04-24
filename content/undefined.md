Do you get an error like this?
> Uncaught TypeError: Cannot read properties of undefined (reading 'name')

This means you have tried to access a property on a value that is `undefined`. A value could be undefined for many reasons.
```javascript
// Using a variable without declaring it:
foo.name;

// Declaring a variable, but not setting a value:
let foo;
foo.name;

// Declaring a variable, but setting it to undefined explicitly...
// or returning undefined from a function or method:
let foo = undefined;
foo.name;
```
You'll have to find the variable on which you try to access the property, and work backwards to find out why it's `undefined`.
Errors come with stack traces right underneath. It will say the exact line and position where the error occurred.
See `/js-error` for the structure of an error message in JavaScript.
