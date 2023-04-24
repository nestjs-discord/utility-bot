How to interpret JavaScript errors:
```text
ReferenceError: num is not defined  <--- Type and message
    at myFunc (variables.js:26:19)  <--- Where it happened
    at inner (other.js:5:9)         <--- Where that was called
    at main (index.js:12:5)         <--- Where THAT was called
    ... etc ...

at myFunc     | variables.js | 26  :  19
   ^^^^^^     | ^^^^^^^^^^^^ | ^^     ^^
function name | file name    | line   character
```
