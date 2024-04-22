import { WebPlugin } from '@capacitor/core';

import type { ClientPlugin, ClientCallback, SendOptions } from './definitions';

function getUrl(model: string) {
  const base = "http://localhost:4242/";
  switch (model) {
    case "dummy":
      return base + "dummy";
    case "openai":
      return base + "openai";
    default:
      return base + "dummy";
  }
}

export class Client extends WebPlugin implements ClientPlugin {


  async send(options: SendOptions, callback: ClientCallback) {
    const url = getUrl(options.model!);
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Connection: 'Keep-Alive'
      },
      body: options.body
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

