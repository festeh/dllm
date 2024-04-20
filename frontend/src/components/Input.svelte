<script lang="ts">
	import { Input } from "svelte-ux";
	import { responseStore } from "../stores/response";
	import { resizeOutputField } from "$lib";

  let input = "";

	async function sendInput() {
		const res = await fetch('http://localhost:4242/openai', {
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
</script>

<div class="p-2 border rounded-xl my-2 mx-4">
	<Input
		placeholder="Type your message here"
		bind:value={input}
		on:keyup={(e) => e.key === 'Enter' && sendInput()}
	/>
</div>
