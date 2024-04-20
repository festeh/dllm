import { resizeOutputField } from "$lib";
import { get } from "svelte/store";
import { responseStore } from "../stores/response";
import { settingsBarStore } from "../stores/settingsBar";

function getUrl() {
  const base = "http://localhost:4242/";
  const settings = get(settingsBarStore);
  switch (settings.selected) {
    case "dummy":
      return base + "dummy";
    case "openai":
      return base + "openai";
    default:
      return base + "dummy";
  }
}

export async function send(input: string) {
  const url = getUrl();
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Connection: 'Keep-Alive'
    },
    body: JSON.stringify({ message: input })
  });
  const reader = res.body!.getReader();
  let done = false;
  responseStore.set('');
  while (!done) {
    const { value, done: d } = await reader!.read();
    done = d;
    if (value != undefined && value.length > 0) {
      const chunk = new TextDecoder().decode(value);
      responseStore.update(chunk)
      resizeOutputField(document.getElementById('response'));
    }
  }
}
