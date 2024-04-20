<script lang="ts">
	import { send } from '$lib/chat';
	import { Button, Input } from 'svelte-ux';

	let input = '';
</script>

<div class="flex p-2 border rounded-xl my-2 mx-4">
	<Input
		placeholder="Type your message here"
		bind:value={input}
		on:keyup={(e) => {
			if (e.ctrlKey && e.key === 'Enter') {
				send(input);
			}
		}}
	/>

	<Button
		icon="m2 21l21-9L2 3v7l15 2l-15 2z"
		on:click={() => send(input)}
		class="p-2"
		target="_blank"
	/>
</div>
