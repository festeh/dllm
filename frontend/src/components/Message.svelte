<script lang="ts">
	import { resizeOutputField } from '$lib';
	import { cls } from 'svelte-ux';
	import { responseStore } from '../stores/response';

	function onInput(e) {
		console.log(e);
		resizeOutputField(e.target);
	}
</script>

<textarea
	id={'response'}
	placeholder={'Waiting...'}
	required
	bind:value={$responseStore}
	on:input={onInput}
	on:focus
	on:blur
	on:keydown
	on:keypress
	class={cls(
		'text-sm border bg-transparent outline-none resize-none',
		'p-2 mx-4 my-2',
		'placeholder-surface-content placeholder-opacity-0 group-focus-within:placeholder-opacity-50'
	)}
/>
