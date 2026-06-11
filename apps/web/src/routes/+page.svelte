<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<main class="scss-proof">
	<h1>SCSS is active</h1>
	<p>If this card has gradient, blur, and rounded corners, SCSS is working.</p>
</main>

<style lang="scss">
	$bg-a: #1b1536;
	$bg-b: #4a2f7a;
	$text: #f5f2ff;

	.scss-proof {
		min-height: 100dvh;
		display: grid;
		place-content: center;
		gap: 0.5rem;
		padding: 2rem;
		color: $text;
		background: radial-gradient(80% 120% at 20% 10%, $bg-b 0%, $bg-a 55%, #09070f 100%);

		h1 {
			margin: 0;
			font-size: clamp(1.6rem, 3vw, 2.4rem);
		}

		p {
			margin: 0;
			opacity: 0.86;
		}

		&::before {
			content: '';
			position: fixed;
			inset: 12vh 10vw;
			z-index: -1;
			border-radius: 24px;
			backdrop-filter: blur(8px);
			box-shadow: 0 24px 80px rgba(0, 0, 0, 0.45);
			background: linear-gradient(140deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.02));
		}
	}
</style>
