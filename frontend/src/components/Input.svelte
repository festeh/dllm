<script lang="ts">
	import { send, stopChat } from '$lib/chat';
	import { Button, Input } from 'svelte-ux';
	import { isChatIdle } from '../stores/chatState';

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

	{#if $isChatIdle}
		<Button icon="m2 21l21-9L2 3v7l15 2l-15 2z" on:click={() => send(input)} class="p-2" />
	{:else}
		<!-- hand-stop -->
		<Button
			icon="M17 7.5a1.5 1.5 0 0 1 3 0V16a6 6 0 0 1-6 6h-2h.208a6 6 0 0 1-5.012-2.7L7 19q-.468-.718-3.286-5.728a1.5 1.5 0 0 1 .536-2.022a1.87 1.87 0 0 1 2.28.28L8 13"
			on:click={stopChat}
		/>
	{/if}
</div>
