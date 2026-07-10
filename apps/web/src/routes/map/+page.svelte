<svelte:head>
	<title>Map</title>
</svelte:head>

<script lang="ts">
	import '$lib/styles/screens/map-screen.scss';
	import MindMap from '$lib/components/MindMap.svelte';
</script>

<div class="map-screen__container">
	<div class="map-screen__map">
		<MindMap />
	</div>
</div>
