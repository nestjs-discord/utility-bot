Windows users who are using TypeScript version 4.9 and up may encounter this problem.
This happens when you're trying to run your application in watch mode, e.g `npm run start:dev` and see an endless loop of the log messages:
```text
XX:XX:XX AM - File change detected. Starting incremental compilation...
XX:XX:XX AM - Found 0 errors. Watching for file changes.
```
When you're using the NestJS CLI to start your application in watch mode it is done by calling `tsc --watch`, and as of version 4.9 of TypeScript, a new strategy for detecting file changes is used which is likely to be the cause of this problem.
In order to fix this problem, you need to add a setting to your tsconfig.json file after the `"compilerOptions"` option as follows:
```json
{
  "watchOptions": {
    "watchFile": "fixedPollingInterval"
  }
}
```
This tells TypeScript to use the polling method for checking for file changes instead of file system events (the new default method), which can cause issues on some machines.
You can read more about the `"watchFile"` option in TypeScript documentation.
