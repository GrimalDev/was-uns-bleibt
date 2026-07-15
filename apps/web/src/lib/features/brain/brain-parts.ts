import type { BrainPartId } from './brain-flow';

export type BrainPart = {
	id: BrainPartId;
	name: string;
};

export const brainParts: BrainPart[] = [
	{ id: 1, name: 'Memories' },
	{ id: 2, name: 'People' },
	{ id: 3, name: 'Places' },
	{ id: 4, name: 'Objects' },
	{ id: 5, name: 'Feelings' }
];
