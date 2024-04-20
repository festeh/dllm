<script lang="ts">
	import { SelectField } from 'svelte-ux';
	import { settingsBarStore } from '../stores/settingsBar';
</script>

<div class="bg-secondary-300 mb-2 p-3">
	<SelectField
		options={$settingsBarStore.options}
    bind:value={$settingsBarStore.selected}
	/>
</div>
