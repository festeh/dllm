import { writable } from "svelte/store";


const store = writable({
  options:
    [
      { value: 'dummy', label: 'Dummy model' },
      { value: 'openai', label: 'OpenAI' }
    ],
  selected: "dummy"
});

export const settingsBarStore = store;
