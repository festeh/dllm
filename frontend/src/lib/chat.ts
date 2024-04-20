import { resizeMessageBox } from "$lib";
import { get } from "svelte/store";
import { settingsBarStore } from "../stores/settingsBar";
import { isChatIdle } from "../stores/chatState";
import { chatHistoryStore } from "../stores/chatHistory";

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
  chatHistoryStore.update((history) => [
    ...history,
    { id: Math.random().toString(), role: "user", content: input }
  ])
  let messages = get(chatHistoryStore).map((message) => {
    return { role: message.role, content: message.content }
  })
  messages = [{ role: "system", content: "Hi!" }, ...messages]
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Connection: 'Keep-Alive'
    },
    body: JSON.stringify({ messages })
  });
  const reader = res.body!.getReader();
  let done = false;
  isChatIdle.set(false);
  const responseId = Math.random().toString()
  chatHistoryStore.update((history) => [
    ...history,
    { id: responseId, role: "assistant", content: '' }
  ])

  while (!done) {
    if (get(isChatIdle)) {
      break
    }
    const { value, done: d } = await reader!.read();
    done = d;
    if (value != undefined && value.length > 0) {
      const chunk = new TextDecoder().decode(value);
      chatHistoryStore.update((history) => {
        const last = history[history.length - 1];
        last.content += chunk;
        return history;
      });
      resizeMessageBox(document.getElementById(responseId));
    }
  }
  isChatIdle.set(true);
}

export async function stopChat() {
  isChatIdle.set(true);
}
