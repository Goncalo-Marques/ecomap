import { get } from "svelte/store";
import { locale } from "./i8n";
import { DateFormats } from "../constants/date";

/**
 * Date used to format times.
 */
const UNIX_TIME = "1970-01-01";

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

/**
 * Formats times to 2-digit hour and 2-digit minutes.
 * @param time Time to be formatted.
 * @returns Formatted time.
 */
export function formatTime(time: string): string {
	return formatDate(`${UNIX_TIME} ${time}`, DateFormats.shortTime);
}
