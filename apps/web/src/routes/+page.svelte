<svelte:head>
	<title>Brain</title>
</svelte:head>

<script lang="ts">
	import FormOverlay from '$lib/components/FormOverlay.svelte';

	const isFormOverlayOpen = true;
</script>

{#if isFormOverlayOpen}
	<FormOverlay />
{/if}
