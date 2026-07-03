<script lang="ts">
	import '$lib/styles/screens/map-screen.scss';
	import '$lib/styles/screens/question-screen.scss';
</script>

<main class="page-root" aria-label="main-view">
	<div class="map-screen" aria-hidden="true"></div>
	<section class="prompt-card" aria-label="question input panel">
		<div class="prompt-copy">
			<p class="prompt-label">Question</p>
			<h1 class="prompt-question">What memory or thought would you like to leave here?</h1>
		</div>

		<label class="answer-field">
			<span class="sr-only">Your answer</span>
			<input
				type="text"
				name="answer"
				autocomplete="off"
				placeholder="For example: I still remember the warmth of your voice — type your own answer..."
			/>
		</label>
	</section>
</main>
