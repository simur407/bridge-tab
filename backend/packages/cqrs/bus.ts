import { logger } from '@packages/logger';
import { Command, CommandHandler } from './command';
import { Query, QueryHandler } from './query';

export interface Bus {
  publish<T extends Command | Query, R = void>(event: T): Promise<R>;
  register<T extends Command>(event: T, handler: CommandHandler<T>): void;
  register<T extends Query>(event: T, handler: QueryHandler<T>): void;
}

class InMemoryBus implements Bus {
  private handlers = new Map<string, CommandHandler<Command> | QueryHandler<Query>>();

  register<T extends Command | Query>(
    event: { new (...args: any[]): T; _name: string },
    handler: CommandHandler<T> | QueryHandler<T>,
  ): void {
    if (this.handlers.has(event._name)) {
      throw new Error(`Handler already registered for ${event._name}`);
    }
    this.handlers.set(event._name, handler);
  }

  async publish<T extends Command | Query, R = void>(event: T): Promise<R> {
    const handler = this.handlers.get(event._name);
    if (!handler) {
      throw new Error(`Handler not registered for ${event._name}`);
    }
    logger.debug(`Publishing event ${event._name}`);
    const result = (await handler.handle(event)) as R;
    logger.debug(`Handled event ${event._name}`);
    return result;
  }
}

export const bus = new InMemoryBus();
