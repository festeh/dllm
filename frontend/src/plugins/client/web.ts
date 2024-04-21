import { WebPlugin } from '@capacitor/core';

import type { ClientPlugin, ClientCallback } from './definitions';

export class Client extends WebPlugin implements ClientPlugin {

  async send(url: string, body: string, callback: ClientCallback): Promise<string> {
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Connection: 'Keep-Alive'
      },
      body
    });
    const reader = res.body!.getReader();
    let done = false;
    while (!done) {
      const { value, done: d } = await reader!.read();
      done = d;
      if (value != undefined && value.length > 0) {
        const chunk = new TextDecoder().decode(value);
        const shouldStop = callback(chunk);
        if (shouldStop) {
          return "stopped"
        }
      }
    }
    return "done"
  }
}

