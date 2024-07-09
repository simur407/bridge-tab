import { bus } from '@packages/cqrs';
import { Request, Response } from 'firebase-functions';
import { CreateTurnamentCommand } from '../../application/commands';

export const createTurnament = async (request: Request, response: Response) => {
  const { id, name } = request.body;

  const command = new CreateTurnamentCommand(id, name);
  await bus.publish(command);
  response.sendStatus(201);
};
