The most common error in Nest is failing to resolve provider dependencies.
The error message usually looks something like this:
```text
Nest can't resolve dependencies of the <provider> (?). Please make sure that the argument <unknown_token> at index [<index>] is available in the <module> context.

Potential solutions:
- Is <module> a valid NestJS module?
- If <unknown_token> is a provider, is it part of the current <module>?
- If <unknown_token> is exported from a separate @Module, is that module imported within <module>?
  @Module({
    imports: [ /* the Module containing <unknown_token> */ ]
  })
```
The most common culprit of the error, is not having the `<provider>` in the module's `providers` array.
Please make sure that the provider is indeed in the `providers` array and following standard NestJS provider practices.

There are a few gotchas, that are common. One is putting a provider in an `imports` array.
If this is the case, the error will have the provider's name where `<module>` should be.
