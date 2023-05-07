Please format your question with Markdown formatting.
It leads to better readability and an easier time to spot problems.
For code blocks, you can wrap your block with three back ticks before and after the block, and after the first three back ticks you can add a language (like `ts`) to add syntax highlighting.
e.g.

\`\`\`ts
@Injectable()
export class MySuperAwesomeService {
  constructor(@Inject('InjectionToken') private readonly dep: SomeDependency) {}

  getRandomNumber(): number {
    return Math.round(Math.random() * 1000);
  }
}
\`\`\`

Becomes

```ts
@Injectable()
export class MySuperAwesomeService {
  constructor(@Inject('InjectionToken') private readonly dep: SomeDependency) {}

  getRandomNumber(): number {
    return Math.round(Math.random() * 1000);
  }
}
```
