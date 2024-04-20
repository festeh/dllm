import { writable } from "svelte/store";

export type ChatMessage = {
  id: string,
  role: "system" | "assistant" | "user",
  content: string
}

export const chatHistoryStore = writable<ChatMessage[]>([])
