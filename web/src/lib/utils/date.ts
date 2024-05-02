import { get } from "svelte/store";
import { locale } from "./i8n";

/**
 * Formats a date into a given format.
 * @param date Date to be formatted.
 * @param format Date format.
 * @returns Formatted date.
 */
export function formatDate(
	date: string,
	format: Intl.DateTimeFormatOptions,
): string {
	const selectedLocale = get(locale);
	const dateFormat = new Intl.DateTimeFormat(selectedLocale, format);

	return dateFormat.format(new Date(date));
}
