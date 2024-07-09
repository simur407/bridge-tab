export abstract class Query {
  static get _name(): string {
    return this.name.slice(0, -'Query'.length);
  }

  get _name(): string {
    return this.constructor.name.slice(0, -'Command'.length);
  }
}

export interface QueryHandler<T extends Query, R = unknown> {
  handle(query: T): Promise<R>;
}
