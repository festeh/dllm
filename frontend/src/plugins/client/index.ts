import { registerPlugin } from '@capacitor/core';

import type { ClientPlugin } from './definitions';

const Client = registerPlugin<ClientPlugin>('Client',
  { "web": () => import('./web').then(m => new m.Client()) }
);

export * from './definitions';
export { Client };
