<script lang="ts">
	import { chatHistoryStore } from '../stores/chatHistory';
	import Input from './Input.svelte';
	import Message from './Message.svelte';
	import SettingsBar from './SettingsBar.svelte';
</script>

<div class="flex flex-col w-full">
	<SettingsBar />
	{#each $chatHistoryStore as message (message.id)}
		<Message {message} />
	{/each}
	<Input />
</div>
