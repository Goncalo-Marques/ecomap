/**
 * Date format.
 * Used as the key names for {@link DateFormats}.
 */
type DateFormat = "shortDateTime" | "shortDate";

/**
 * Date formats.
 */
export const DateFormats: Record<DateFormat, Intl.DateTimeFormatOptions> = {
	shortDateTime: {
		day: "2-digit",
		month: "2-digit",
		year: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	},
	shortDate: {
		day: "2-digit",
		month: "2-digit",
		year: "numeric",
	},
};
