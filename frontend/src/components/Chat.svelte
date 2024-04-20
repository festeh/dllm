<script lang="ts">
	import { Button } from 'svelte-ux';
	import Input from './Input.svelte';
	import Message from './Message.svelte';
	import SettingsBar from './SettingsBar.svelte';
	import { chatHistoryStore } from '../stores/chatHistory';
	import { clearChat } from '$lib/chat';
</script>

<div class="flex flex-col w-full">
	<SettingsBar />
	{#each $chatHistoryStore as message (message.id)}
		<Message {message} />
	{/each}
	<Input />
	<Button class="btn border mx-4 p-2" on:click={clearChat}>Clear</Button>
</div>
