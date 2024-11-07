import { register } from 'tsconfig-paths';
register({
  baseUrl: './jdkfajfkfja',
  paths: {
    '@modules/*': ['modules/*'],
    '@packages/*': ['packages/*'],
    '@bridge-tab/*': ['*'],
  },
});
import { TurnamentManagementModule } from '@modules/turnament-management';
import { Module } from '@packages/modules';
import * as functions from 'firebase-functions';
import * as logger from 'firebase-functions/logger';

const modules: Module[] = [TurnamentManagementModule.register()];

export const api = modules.reduce((exports, module) => {
  const moduleControllers = Object.entries(module.controllers).reduce<Record<string, functions.HttpsFunction>>(
    (exports, [name, controller]) => {
      exports[name] = functions.region('europe-west1').https.onRequest(controller);
      return exports;
    },
    {},
  );

  return {
    ...exports,
    [module.name]: moduleControllers,
  };
}, {});
