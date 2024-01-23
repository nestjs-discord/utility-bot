Nest applications handle requests and produce responses in a sequence we refer to as the **request lifecycle**.
With the use of middleware, pipes, guards, and interceptors, it can be challenging to track down where a particular piece of code executes during the request lifecycle, especially as global, controller level, and route level components come into play.
### In general, the request lifecycle looks like the following:

1. Incoming request
2. [Middleware](<https://docs.nestjs.com/middleware>)
   - Globally bound middleware
   - Module bound middleware
3. [Guards](<https://docs.nestjs.com/guards>)
   - Global guards
   - Controller guards
   - Route guards
4. [Interceptors](<https://docs.nestjs.com/interceptors>) _(pre-controller)_
   - Global interceptors
   - Controller interceptors
   - Route interceptors
5. [Pipes](<https://docs.nestjs.com/pipes>)
   - Global pipes
   - Controller pipes
   - Route pipes
   - Route parameter pipes
6. [Controller](<https://docs.nestjs.com/controllers>) _(method handler)_
7. [Service](<https://docs.nestjs.com/providers>) _(if exists)_
8. [Interceptors](<https://docs.nestjs.com/interceptors>) _(post-request)_
   - Route interceptor
   - Controller interceptor
   - Global interceptor
9. [Exception filters](<https://docs.nestjs.com/exception-filters>)
   - Route
   - Controller
   - Global
10. Server response
