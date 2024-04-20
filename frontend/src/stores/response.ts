import { writable } from "svelte/store";


const store = writable("");

export const responseStore = {
  subscribe: store.subscribe,
  set: (message: string) => {
    store.set(message);
  },
  update: (message: string) => {
    store.update((state) => {
      return state + message;
    });
  },
  clear: () => {
    store.set("");
  }
};
