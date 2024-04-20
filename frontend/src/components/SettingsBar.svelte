<script lang="ts">
	import { SelectField } from 'svelte-ux';
	import { settingsBarStore } from '../stores/settingsBar';
	import { clearChat } from '$lib/chat';
	import { isChatIdle } from '../stores/chatState';
</script>

<div class="bg-secondary-300 mb-2 p-3">
	<SelectField
		options={$settingsBarStore.options}
		bind:value={$settingsBarStore.selected}
		on:change={(e) => {
			isChatIdle.set(true);
			clearChat();
		}}
	/>
</div>
