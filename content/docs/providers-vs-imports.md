Don't list a service in `providers` array of a module unless the module really *provides* the service.
If your `ServiceA` depends on a `ServiceB` from different module `ModuleB`, add `ServiceB` to `ModuleB`s `exports` array, 
then add `ModuleB` to `imports` array in your module.
```ts
@Module({
  providers: [ServiceB],
  exports: [
    ServiceB // <-- export the "ServiceB" provider by adding it to the module's exports array
  ], 
})
export class ModuleB {}
```
```ts
@Module({
  imports: [
    ModuleB // <-- modules that wish to inject the "ServiceB" will need to import the "ModuleB" in their imports array
  ],  
  providers: [ServiceA],
})
export class ModuleA {}
```
```ts
import { ServiceB } from '../module-b.module'

@Injectable()
export class ServiceA {
  constructor(private readonly serviceB: ServiceB) {}
  
  // async foo() {
  //   this.serviceB.bar();
  // }
}
```
