export type FormOption = {
	titre: string;
	placeholder: string;
};

export type FormPart = {
	id: number;
	name: string;
	options: FormOption[];
};

export async function loadFormDefinition(): Promise<FormPart[]> {
	const response = await fetch('/data/form.json');
	if (!response.ok) return [];

	return (await response.json()) as FormPart[];
}

export function selectRandomOption(formDefinition: FormPart[], partId: number): FormOption | null {
	const part = formDefinition.find((candidate) => candidate.id === partId);
	if (!part?.options.length) return null;

	return part.options[Math.floor(Math.random() * part.options.length)] ?? null;
}
