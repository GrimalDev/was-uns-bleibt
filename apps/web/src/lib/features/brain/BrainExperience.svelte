<script lang="ts">
	import AnswerAnimation from './AnswerAnimation.svelte';
	import BrainView from './BrainView.svelte';
	import PartQuestion from './PartQuestion.svelte';
	import SaveAnimation from './SaveAnimation.svelte';
	import type { BrainFlowState, BrainPartId } from './brain-flow';

	let state: BrainFlowState = 'brain';
	let selectedPartId: BrainPartId | null = null;
	let answer = '';

	function selectPart(partId: BrainPartId) {
		selectedPartId = partId;
		state = 'transition-to-part';
	}

	function finishTransition() {
		state = 'question';
	}

	function submitAnswer(nextAnswer: string) {
		answer = nextAnswer;
		state = 'save-transition';
	}

	function finishSaveTransition() {
		selectedPartId = null;
		answer = '';
		state = 'brain';
	}
</script>

<main class="brain-experience" aria-label="Brain experience">
	{#if state === 'brain'}
		<BrainView onSelectPart={selectPart} />
	{:else if state === 'transition-to-part'}
		<AnswerAnimation onComplete={finishTransition} />
	{:else if state === 'question' && selectedPartId !== null}
		<PartQuestion partId={selectedPartId} onAnswer={submitAnswer} />
	{:else if state === 'save-transition'}
		<SaveAnimation onComplete={finishSaveTransition} />
	{/if}
</main>


<style lang="scss">
	.brain-experience {
		position: fixed;
		inset: 0;
		overflow: hidden;
	}
</style>
