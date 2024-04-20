<script lang="ts">
	import { send } from '$lib/chat';
	import { Input } from 'svelte-ux';

	let input = '';
</script>

<div class="p-2 border rounded-xl my-2 mx-4">
	<Input
		placeholder="Type your message here"
		bind:value={input}
		on:keyup={(e) => e.key === 'Enter' && send(input)}
	/>
</div>
