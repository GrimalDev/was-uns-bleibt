<script lang="ts">
	import { language, availableLanguages } from '$lib/stores/language';

	const flags = import.meta.glob('../assets/flags/*.png', {
		eager: true,
		import: 'default',
		query: '?url'
	}) as Record<string, string>;

	function getFlag(languageCode: string): string {
		return flags[`../assets/flags/${languageCode}.png`] ?? '';
	}
</script>

<nav class="language-selector" aria-label="Language selection">
	{#each availableLanguages as languageCode}
		<button
			type="button"
			class:active={$language === languageCode}
			aria-label={`Switch language to ${languageCode}`}
			aria-pressed={$language === languageCode}
			on:click={() => language.set(languageCode)}
		>
			<img src={getFlag(languageCode)} alt="" aria-hidden="true" />
		</button>
	{/each}
</nav>

<style lang="scss">
	.language-selector {
		position: absolute;
		top: 0.6em;
		right: 20vw;
		z-index: 1000;
		display: flex;
		gap: 0.5rem;
	}

	button {
		display: block;
		padding: 0;
		border: 0;
		outline: 0;
		background: transparent;
		cursor: pointer;

		&:hover,
		&:focus-visible {
			transform: scale(1.5);
		}
	}

	.language-selector:not(:has(button:hover)) button.active {
		transform: scale(1.5);
	}

	img {
		display: block;
		width: 2.5em;
		height: 2em;
		object-fit: cover;
	}
</style>
