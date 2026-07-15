<script lang="ts">
	import { onMount } from 'svelte';

	export let onComplete: () => void;

	const TRANSITION_DURATION_MS = 1500;
	let timeout: ReturnType<typeof setTimeout>;

	onMount(() => () => clearTimeout(timeout));

	function completeTransition() {
		clearTimeout(timeout);
		timeout = setTimeout(onComplete, TRANSITION_DURATION_MS);
	}
</script>

<section
	class="answer-animation-placeholder"
		aria-live="polite"
	>
		<img
			class="answer-animation-placeholder__gif"
			src="/style/animations/brain_transition.gif"
			alt=""
			on:load={completeTransition}
		/>
	</section>

<style lang="scss">
	.answer-animation-placeholder {
		display: grid;
		place-items: center;
		width: 100%;
		height: 100%;
		position: relative;
		color: white;

		&__gif {
			position: absolute;
			inset: 0;
			width: 100%;
			height: 100%;
			object-fit: contain;
		}
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>
