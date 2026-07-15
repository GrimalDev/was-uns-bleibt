<script lang="ts">
	import { onMount } from 'svelte';
	import '$lib/styles/screens/question-screen.scss';
	import type { BrainPartId } from './brain-flow';
	import { formDefinition, loadFormDefinition, selectRandomOption } from './form-data';

	const MAX_ANSWER_LENGTH = 40;

	export let partId: BrainPartId;
	export let onAnswer: (answer: string) => void;

	let answer = '';
	let selectedPartId: BrainPartId | null = null;
	let selectedOptionId: number | null = null;
	let questionTitle = 'What memory or thought would you like to leave here?';
	let questionPlaceholder =
		'For example: I still remember the warmth of your voice — type your own answer...';

	$: if ($formDefinition.length > 0 && selectedPartId !== partId) {
		selectedPartId = partId;
		selectedOptionId = selectRandomOption($formDefinition, partId)?.id ?? null;
	}

	$: selectedOption = $formDefinition
		.find((part) => part.id === partId)
		?.options.find((option) => option.id === selectedOptionId);
	$: questionTitle = selectedOption?.titre || 'What memory or thought would you like to leave here?';
	$: questionPlaceholder =
		selectedOption?.placeholder ||
		'For example: I still remember the warmth of your voice — type your own answer...';

	$: questionSegments = questionTitle.match(/\*[^*]+\*|[^*]+/g) ?? [questionTitle];
	$: questionAccent = `var(--color-brain-${partId})`;

	onMount(async () => {
		await loadFormDefinition();
	});

	function submitAnswer() {
		const trimmedAnswer = answer.trim();
		if (!trimmedAnswer) return;

		onAnswer(trimmedAnswer);
	}
</script>

<div class="form-page" aria-label="question view" data-form-sections={$formDefinition.length}>
	<div class="prompt-layer fade-from-black" style={`--question-accent: ${questionAccent};`}>
		<form class="prompt-card" aria-label="question input panel" on:submit|preventDefault={submitAnswer}>
			<div class="prompt-copy">
				<p class="prompt-label">Question</p>
				<h1 class="prompt-question">
					{#each questionSegments as segment, index (index)}
						{#if segment.startsWith('*') && segment.endsWith('*')}
							<span class="prompt-question__highlight">{segment.slice(1, -1)}</span>
						{:else}
							{segment}
						{/if}
					{/each}
				</h1>
			</div>

			<div class="answer-field">
				<label>
					<span class="sr-only">Your answer</span>
					<input
						bind:value={answer}
						type="text"
						name="answer"
						autocomplete="off"
						placeholder={questionPlaceholder}
						maxlength={MAX_ANSWER_LENGTH}
						aria-describedby="answer-character-count"
						required
					/>
				</label>
				<p
					id="answer-character-count"
					class:maxed={answer.length >= MAX_ANSWER_LENGTH}
					class="character-count"
					aria-live="polite"
				>
					Max: {answer.length} / {MAX_ANSWER_LENGTH}
				</p>
			</div>
		</form>
	</div>
</div>
