import { writable } from 'svelte/store';

export const availableLanguages = ['FR', 'DE'] as const;
export type FormLanguage = (typeof availableLanguages)[number];

export const language = writable<FormLanguage>('FR');
