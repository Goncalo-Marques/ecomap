import { derived, writable } from "svelte/store";
import en from "../../locales/en.json";
import pt from "../../locales/pt.json";
import schema from "../../locales/schema.json";

/**
 * Identifiers of the locale schema.
 */
type LocaleTextID = keyof (typeof schema)["properties"];

/**
 * Map of the available locales with their respective configuration.
 */
const locales: Record<string, Record<LocaleTextID, string>> = {
	en,
	pt,
};

/**
 * Localizes text to a given locale.
 * @param locale Locale to localize a given text.
 * @param id Identifier of the text to be localized.
 * @param vars Map used to replace placeholder values inside a given text.
 * @returns Localized text.
 */
function translate(
	locale: string,
	id: LocaleTextID,
	vars: Record<string, unknown>,
) {
	const localeMap = locales[locale];
	let text = localeMap[id];

	// Replace placeholder variables inside text.
	for (const key of Object.keys(vars)) {
		text = text.replace(/{{.*}}/g, String(vars[key]));
	}

	return text;
}

/**
 * Locale store for reading and updating the locale of the application.
 */
export const locale = writable("en");

/**
 * Derived store which translates a given text to the locale configured in the application.
 */
export const t = derived(locale, function ($locale) {
	/**
	 * Translates a given text to the locale configured in the application.
	 * @param id Identifier of the text to be localized.
	 * @param vars Map used to replaced placeholder values inside a given text.
	 */
	return function (id: LocaleTextID, vars: Record<string, unknown> = {}) {
		return translate($locale, id, vars);
	};
});
