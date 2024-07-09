import { Command, CommandHandler } from '@packages/cqrs';
import { Turnament, TurnamentRepository } from '../../domain/turnament';

export class CreateTurnamentCommand extends Command {
  constructor(readonly id: string, readonly name: string) {
    super();
  }
}

export class CreateTurnamentHandler implements CommandHandler<CreateTurnamentCommand> {
  constructor(private readonly turnamentRepository: TurnamentRepository) {}
  async handle(command: CreateTurnamentCommand): Promise<void> {
    const turnament = Turnament.create({
      id: command.id,
      name: command.name,
    });

    await this.turnamentRepository.save(turnament);
  }
}
