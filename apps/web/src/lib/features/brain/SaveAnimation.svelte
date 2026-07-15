<script lang="ts">
	import { onMount } from 'svelte';

	export let onComplete: () => void;

	const SAVE_TRANSITION_DURATION_MS = 2000;
	let timeout: ReturnType<typeof setTimeout>;

	onMount(() => () => clearTimeout(timeout));

	function completeTransition(): void {
		clearTimeout(timeout);
		timeout = setTimeout(onComplete, SAVE_TRANSITION_DURATION_MS);
	}
</script>

<section class="save-animation" aria-live="polite">
	<img
		class="save-animation__gif"
		src="/style/animations/brain_save_transition.gif"
		alt=""
		on:load={completeTransition}
	/>
</section>

<style lang="scss">
	.save-animation {
		position: relative;
		display: grid;
		place-items: center;
		width: 100%;
		height: 100%;
    background: url('/style/backgrounds/fill_background.png') no-repeat center;
    background-size: cover;
		color: white;

		&__gif {
			position: absolute;
			inset: 0;
			width: 100%;
			height: 100%;
			object-fit: contain;
		}
	}
</style>
