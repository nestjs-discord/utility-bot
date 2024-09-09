When you use the `ConfigModule` from the `@nestjs/config` package, and you need to access the `process.env.FOO` variable in another module's registration function or constructors, you might encounter issues where `process.env.FOO` is undefined, even though it has the correct value when the application is fully bootstrapped.
This occurs because `process.env` is not populated when you are accessing it. In such a situation, there are generally two solutions:

## **Solution 1.** ✅
If the module provides an asynchronous way to set its options using the `forRootAsync()` or `registerAsync()` method, it's better to use this along with the `ConfigService` instead of directly reading `process.env`.
```ts
import { JwtModule } from '@nestjs/jwt';

@Module({
  imports: [
    JwtModule.registerAsync({
      imports: [ConfigModule],
      useFactory: async (configService: ConfigService) => ({
        secretOrPrivateKey: configService.get<string>('jwt.secret'),
      }),
      inject: [ConfigService],
    }),
  ]
})
export class FooModule {}
```
## **Solution 2.** ✅
In the `main.ts`, the **very first** lines must be
```ts
import 'dotenv/config';
// or, the line below:
// import { config } from 'dotenv'; config();
```
It should also be __before__ **all other imports**. This will populate `process.env` for you.

### Side note: ⚠️
When using the `@nestjs/jwt` package, avoid registering the `JwtService` in a `providers` array. It's not necessary because of the import of the `JwtModule`. Only import the `JwtModule` in the `imports` array and remove the `JwtService` from the `providers` array when defining the modules.
```diff
import {
+  JwtModule,
-  JwtService
} from '@nestjs/jwt';

@Module({
  imports: [
+   JwtModule.registerAsync({
+     // ...
+   }),
  ],
  providers: [
-   JwtService,
  ],
})
export class FooModule {}
```
