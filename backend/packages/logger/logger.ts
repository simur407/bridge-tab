import * as firebaseLogger from 'firebase-functions/logger';

export interface Logger {
  debug(message: string, data?: object): void;
  info(message: string, data?: object): void;
  error(message: string, data?: object): void;
}

export class FirebaseLogger implements Logger {
  debug(message: string, data?: object): void {
    firebaseLogger.debug(message, data);
  }
  info(message: string, data?: object): void {
    firebaseLogger.info(message, data);
  }
  error(message: string, data?: object): void {
    firebaseLogger.error(message, data);
  }
}

export const logger = new FirebaseLogger();
