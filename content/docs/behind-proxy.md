When your NestJS application is behind a reverse proxy, you may need to enable the "trust proxy" option in Express or Fastify, depending on your server.

In Express, you can enable the `trust proxy` option using the `app.set()` method.
```typescript
async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.set('trust proxy', true); // <----

  await app.listen(3000);
}
```
If you're using Fastify, you can enable the `trustProxy` option in `FastifyAdapter`.
```typescript
import { FastifyAdapter, NestFastifyApplication } from '@nestjs/platform-fastify';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,

    new FastifyAdapter({
      trustProxy: true, // <----
    }),
  );
}
```
> By enabling the "trust proxy" option, Express/Fastify will trust the `X-Forwarded-For` header and consider it the client's IP address.

Suppose you use the `@nestjs/throttler` module to rate-limit users and protect applications from brute-force attacks.
In that case, besides enabling the "trust proxy" option, you should override the `getTracker()` method to pull the value from the header rather than from `req.ip`.
```ts
// throttler-behind-proxy.guard.ts
import { ThrottlerGuard } from '@nestjs/throttler';
import { Injectable } from '@nestjs/common';

@Injectable()
export class ThrottlerBehindProxyGuard extends ThrottlerGuard {
  protected getTracker(req: Record<string, any>): string {
    return req.ips.length ? req.ips[0] : req.ip; // individualize IP extraction to meet your own needs
  }
}

// app.controller.ts
import { ThrottlerBehindProxyGuard } from './throttler-behind-proxy.guard';

@UseGuards(ThrottlerBehindProxyGuard)
```
