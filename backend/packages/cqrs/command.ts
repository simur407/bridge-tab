export abstract class Command {
  static get _name(): string {
    return this.name.slice(0, -'Command'.length);
  }

  get _name(): string {
    return this.constructor.name.slice(0, -'Command'.length);
  }
}

export interface CommandHandler<T extends Command> {
  handle(command: T): Promise<void>;
}
