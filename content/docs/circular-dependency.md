Occasionally you'll find it difficult to avoid circular dependencies in your application.
You'll need to take some steps to help Nest resolve these. Errors that arise from circular dependencies look like this:
```text
Nest cannot create the <module> instance.
The module at index [<index>] of the <module> "imports" array is undefined.

Potential causes:
- A circular dependency between modules. Use forwardRef() to avoid it. Read more: https://docs.nestjs.com/fundamentals/circular-dependency
- The module at index [<index>] is of type "undefined". Check your import statements and the type of the module.

Scope [<module_import_chain>]
# example chain AppModule -> FooModule
```
A circular dependency occurs when two classes depend on each other.
For example, class A needs class B, and class B also needs class A.
Circular dependencies can arise in Nest between modules and between providers.

While circular dependencies should be avoided where possible, you can't always do so.
In such cases, Nest enables resolving circular dependencies between providers in two ways.
In this chapter, we describe using **forward referencing** as one technique, and using the **ModuleRef** class to retrieve a provider instance from the DI container as another.
