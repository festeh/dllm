import { resizeMessageBox } from "$lib";
import { get } from "svelte/store";
import { settingsBarStore } from "../stores/settingsBar";
import { isChatIdle } from "../stores/chatState";
import { chatHistoryStore } from "../stores/chatHistory";

import { Client } from '../plugins/client/'

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

function produceCallback(responseId: string) {
  return function clientCallback(chunk: string)  {
    chatHistoryStore.update((history) => {
      const last = history[history.length - 1];
      last.content += chunk;
      return history;
    });
    resizeMessageBox(document.getElementById(responseId));
    return get(isChatIdle)
  }
}

export async function send(input: string) {
  chatHistoryStore.update((history) => [
    ...history,
    { id: Math.random().toString(), role: "user", content: input }
  ])
  let messages = get(chatHistoryStore).map((message) => {
    return { role: message.role, content: message.content }
  })
  messages = [{ role: "system", content: "Hi!" }, ...messages]
  const body = JSON.stringify({ messages });
  const responseId = Math.random().toString()
  chatHistoryStore.update((history) => [
    ...history,
    { id: responseId, role: "assistant", content: '' }
  ])
  isChatIdle.set(false);
  const cb = produceCallback(responseId);
  const result = await Client.send(getUrl(), body, cb);
  isChatIdle.set(true);
}

export async function stopChat() {
  isChatIdle.set(true);
}

export function clearChat() {
  chatHistoryStore.set([]);
}
