import { writable } from "svelte/store";

type ChatMessage = {
  id: string,
  role: "assistant" | "user",
  message: string
}

export const chatHistoryStore = writable<ChatMessage[]>([])
