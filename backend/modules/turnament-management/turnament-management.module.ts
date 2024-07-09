import { bus } from '@packages/cqrs';
import { Module } from '@packages/modules';
import { CreateTurnamentCommand, CreateTurnamentHandler } from './application/commands';
import { createTurnament } from './infrastructure/api';
import { InMemoryTurnamentRepository } from './domain/turnament';

export class TurnamentManagementModule implements Module {
  name = 'turnament-management';
  controllers = {
    'create-turnament': createTurnament,
  };

  static register() {
    return new TurnamentManagementModule();
  }

  constructor() {
    const turnamentRepository = new InMemoryTurnamentRepository();
    bus.register(CreateTurnamentCommand, new CreateTurnamentHandler(turnamentRepository));
  }
}
