import { writable } from 'svelte/store';

export type FormLanguage = 'FR' | 'DE';

export const language = writable<FormLanguage>('FR');
