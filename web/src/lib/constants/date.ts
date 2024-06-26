/**
 * Date format.
 * Used as the key names for {@link DateFormats}.
 */
type DateFormat = "shortDateTime" | "shortDate" | "shortTime" | "shortTime24H";

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
	shortTime: {
		hour: "2-digit",
		minute: "2-digit",
	},
	shortTime24H: {
		hour: "2-digit",
		minute: "2-digit",
		hourCycle: "h24",
	},
};
