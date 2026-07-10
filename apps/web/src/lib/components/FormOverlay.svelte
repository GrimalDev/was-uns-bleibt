<script lang="ts">
	import formDefinition from '$lib/data/form.json';
	import '$lib/styles/screens/map-screen.scss';
	import '$lib/styles/screens/question-screen.scss';

	const selectedPartId = 1;
	const selectedPart = formDefinition.find((part) => part.id === selectedPartId);
	const selectedOption =
		selectedPart && selectedPart.options.length > 0
			? selectedPart.options[Math.floor(Math.random() * selectedPart.options.length)]
			: null;
	const questionTitle = selectedOption?.titre || 'What memory or thought would you like to leave here?';
	const questionPlaceholder =
		selectedOption?.placeholder ||
		'For example: I still remember the warmth of your voice — type your own answer...';
	const questionSegments = questionTitle.match(/\*[^*]+\*|[^*]+/g) ?? [questionTitle];
	const questionAccent = `var(--color-brain-${selectedPartId})`;
</script>

<div class="form-page" aria-label="main-view" data-form-sections={formDefinition.length}>
	<div class="map-screen" aria-hidden="true"></div>
	<div class="prompt-layer" style={`--question-accent: ${questionAccent};`}>
		<section class="prompt-card" aria-label="question input panel">
			<div class="prompt-copy">
				<p class="prompt-label">Question</p>
				<h1 class="prompt-question">
					{#each questionSegments as segment}
						{#if segment.startsWith('*') && segment.endsWith('*')}
							<span class="prompt-question__highlight">{segment.slice(1, -1)}</span>
						{:else}
							{segment}
						{/if}
					{/each}
				</h1>
			</div>

			<label class="answer-field">
				<span class="sr-only">Your answer</span>
				<input
					type="text"
					name="answer"
					autocomplete="off"
					placeholder={questionPlaceholder}
				/>
			</label>
		</section>
	</div>
</div>
