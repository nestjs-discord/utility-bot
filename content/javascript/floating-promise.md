A "floating" Promise is a Promise that is created without any code set up to handle any errors it might throw. Floating Promises can cause several issues, such as improperly sequenced operations and ignored Promise rejections.

In the example below, a Promise is created and executed. However, since there is no `await` keyword to wait for its resolution, the program continues to run without waiting for the Promise to complete.
```ts
async function someAsyncFunction() {
  throw new Error('Whoops! 😥');
}

try {
  someAsyncFunction(); // no await
} catch (e) {
  // It won't catch the error
}

// Uncaught (in promise) Error: Whoops! 😥
// at someAsyncFunction (<anonymous>:2:9)
// at <anonymous>:6:3
```
The Promise returned by `someAsyncFunction()` starts executing, but the code execution continues immediately without waiting for the Promise to resolve.
As a result, any error or thrown exception inside `someAsyncFunction()` will not be caught by the `try/catch` block because the Promise is not being awaited. This can lead to unexpected behavior or unhandled rejections, as the code does not handle the asynchronous operation correctly.

To ensure that the Promise is awaited and its result is properly handled, you should use the `await` keyword when invoking an asynchronous function:
```ts
try {
  await someAsyncFunction(); // properly awaited
} catch (e) {
  // It will catch the error
  console.log(e.message); // Whoops! 😥
}
```
