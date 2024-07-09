import { Request, Response } from 'firebase-functions';

export type Handler = (req: Request, response: Response) => void | Promise<void>;

export abstract class Module {
  abstract name: string;
  abstract controllers: Record<string, Handler>;
}
