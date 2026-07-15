import { writable } from 'svelte/store';
import { language, type FormLanguage } from '$lib/stores/language';

export type { FormLanguage } from '$lib/stores/language';

type LocalizedText = Record<FormLanguage, string>;

type RawFormOption = {
	id: number;
	titre: LocalizedText;
	placeholder: LocalizedText;
};

type RawFormPart = {
	id: number;
	name: LocalizedText;
	options: RawFormOption[];
};

export type FormOption = {
	id: number;
	titre: string;
	placeholder: string;
};

export type FormPart = {
	id: number;
	name: string;
	options: FormOption[];
};

export const formDefinition = writable<FormPart[]>([]);

let rawFormDefinition: RawFormPart[] = [];
let loadPromise: Promise<void> | null = null;
let currentLanguage: FormLanguage = 'FR';

function resolveFormDefinition(form: RawFormPart[], currentLanguage: FormLanguage): FormPart[] {
	return form.map((part) => ({
		id: part.id,
		name: part.name[currentLanguage],
		options: part.options.map((option) => ({
			id: option.id,
			titre: option.titre[currentLanguage],
			placeholder: option.placeholder[currentLanguage]
		}))
	}));
}

language.subscribe((nextLanguage) => {
	currentLanguage = nextLanguage;
	formDefinition.set(resolveFormDefinition(rawFormDefinition, nextLanguage));
});

export async function loadFormDefinition(): Promise<void> {
	if (loadPromise) return loadPromise;

	loadPromise = (async () => {
		try {
			const response = await fetch('/data/form.json');
			rawFormDefinition = response.ok ? ((await response.json()) as RawFormPart[]) : [];
		} catch {
			rawFormDefinition = [];
		}

		formDefinition.set(resolveFormDefinition(rawFormDefinition, currentLanguage));
	})();

	return loadPromise;
}

export function selectRandomOption(formDefinition: FormPart[], partId: number): FormOption | null {
	const part = formDefinition.find((candidate) => candidate.id === partId);
	if (!part?.options.length) return null;

	return part.options[Math.floor(Math.random() * part.options.length)] ?? null;
}
