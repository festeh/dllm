<script lang="ts">
	import { Input, cls } from 'svelte-ux';
	let response = '';
	let input = '';

	function resize(target: HTMLElement) {
		if (target.clientHeight >= target.scrollHeight && target.clientHeight < 50) {
			return;
		}
		target.style.height = '10px';
		target.style.height = +target.scrollHeight + 'px';
	}

	function onInput(e) {
		console.log(e);
		resize(e.target);
	}

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
		response = '';
		while (!done) {
			const { value, done: d } = await reader!.read();
			done = d;
			if (value != undefined && value.length > 0) {
				const chunk = new TextDecoder().decode(value);
				response += chunk;
				resize(document.getElementById('response'));
			}
		}
	}
</script>

<main class="p-2">
	<h1 class="text-lg font-semibold text-center">Chat</h1>
	<div class="flex flex-col w-full">
		<div class="p-2 border rounded-xl my-2 mx-4">
			<Input
				placeholder="Type your message here"
				bind:value={input}
				on:keyup={(e) => e.key === 'Enter' && sendInput()}
			/>
		</div>
		<textarea
			id={'response'}
			placeholder={'Waiting...'}
			required
			bind:value={response}
			on:input={onInput}
			on:focus
			on:blur
			on:keydown
			on:keypress
			class={cls(
				'text-sm border bg-transparent outline-none resize-none',
				'p-2 mx-4 my-2',
				'placeholder-surface-content placeholder-opacity-0 group-focus-within:placeholder-opacity-50',
				settingsClasses.input
			)}
		/>
	</div>
</main>
