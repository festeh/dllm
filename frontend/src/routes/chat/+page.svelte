<script lang="ts">
	import { Input } from 'svelte-ux';
	let response = '';
	let input = '';

	async function sendInput() {
		console.log('Sending', input);
		const res = await fetch('http://localhost:4242/openai', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Connection: 'Keep-Alive'
			},
			body: JSON.stringify({ message: input })
		});
		const reader = res.body!.getReader();
		console.log(reader);
		let done = false;
		response = '';
		while (!done) {
			console.log('reading');
			const { value, done: d } = await reader!.read();
			done = d;
			if (value) {
				response += new TextDecoder().decode(value);
			}
		}
	}
</script>

<main class="p-2">
	<h1 class="text-lg font-semibold text-center">Chat</h1>
	<div class="p-2 border rounded-xl my-2 mx-4">
		<Input
			placeholder="Type your message here"
			bind:value={input}
			on:keyup={(e) => e.key === 'Enter' && sendInput()}
		/>
	</div>
	<div class="border border-primary-400 h-52">
		{response}
	</div>
</main>
