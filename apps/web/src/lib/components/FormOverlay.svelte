<script lang="ts">
	import { onMount } from 'svelte';
	import '$lib/styles/screens/map-screen.scss';
	import '$lib/styles/screens/question-screen.scss';

	type FormOption = {
		titre: string;
		placeholder: string;
	};

	type FormPart = {
		id: number;
		options: FormOption[];
	};

	let formDefinition: FormPart[] = [];
	const selectedPartId = 1;
	let questionTitle = 'What memory or thought would you like to leave here?';
	let questionPlaceholder =
		'For example: I still remember the warmth of your voice — type your own answer...';
	let questionSegments = questionTitle.match(/\*[^*]+\*|[^*]+/g) ?? [questionTitle];
	const questionAccent = `var(--color-brain-${selectedPartId})`;

	onMount(async () => {
		const response = await fetch('/data/form.json');
		if (!response.ok) return;

		formDefinition = (await response.json()) as FormPart[];
		const selectedPart = formDefinition.find((part) => part.id === selectedPartId);
		const selectedOption = selectedPart?.options.length
			? selectedPart.options[Math.floor(Math.random() * selectedPart.options.length)]
			: null;

		questionTitle = selectedOption?.titre || questionTitle;
		questionPlaceholder = selectedOption?.placeholder || questionPlaceholder;
		questionSegments = questionTitle.match(/\*[^*]+\*|[^*]+/g) ?? [questionTitle];
	});
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
