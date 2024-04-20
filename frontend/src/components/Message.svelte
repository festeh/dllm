<script lang="ts">
	import { resizeMessageBox } from '$lib';
	import { cls } from 'svelte-ux';
	import { onMount } from 'svelte';
	export let message;

	onMount(() => {
		const elem = document.getElementById(message.id);
		if (elem) {
			resizeMessageBox(elem);
		}
	});

	function onInput(e) {
		console.log(e);
		resizeMessageBox(e.target);
	}
</script>

<textarea
	id={message.id}
	bind:value={message.content}
	on:input={onInput}
	on:focus
	on:blur
	on:keydown
	on:keypress
	class={cls(
		'border bg-transparent outline-none resize-none',
		'p-2 mx-4 my-2',
		'placeholder-surface-content placeholder-opacity-0 group-focus-within:placeholder-opacity-50',
		'overflow-hidden'
	)}
/>
