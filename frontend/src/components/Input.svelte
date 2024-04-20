<script lang="ts">
	import { send, stopChat } from '$lib/chat';
	import { Button } from 'svelte-ux';
	import { isChatIdle } from '../stores/chatState';
	import { mdiStopCircleOutline } from '@mdi/js';
	import { resizeMessageBox } from '$lib';

	let input = '';
</script>

<div class="flex p-2 border rounded-xl my-2 mx-4">
	<textarea
		class="flex-1 bg-transparent outline-none resize-none p-2 overflow-hidden"
		placeholder="Type your message here"
		bind:value={input}
		on:keyup={(e) => {
			resizeMessageBox(e.target);
			if (e.ctrlKey && e.key === 'Enter') {
				send(input);
				input = '';
			}
		}}
	/>

	{#if $isChatIdle}
		<Button
			icon="m2 21l21-9L2 3v7l15 2l-15 2z"
			on:click={() => {
				send(input);
				input = '';
			}}
			class="p-2"
		/>
	{:else}
		<Button icon={mdiStopCircleOutline} on:click={stopChat} />
	{/if}
</div>
