Don't list a service in `providers` array of a module unless the module really *provides* the service.
If your `ServiceA` depends on a `ServiceB` from different module `ModuleB`, add `ServiceB` to `ModuleB`s `exports` array, 
then add `ModuleB` to `imports` array in your module.
