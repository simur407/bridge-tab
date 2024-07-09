import { TurnamentManagementModule } from '@modules/turnament-management';
import { Module } from '@packages/modules';
import * as functions from 'firebase-functions';
import * as logger from 'firebase-functions/logger';

const modules: Module[] = [TurnamentManagementModule.register()];

// const functions = modules.map(module => module.controllers());

export const api = functions.region('europe-west1').https.onRequest((request, response) => {
  logger.info('Hello logs!', { structuredData: true });
  response.send('Hello from Firebase!');
});
