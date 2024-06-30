import { derived, get, writable, type Writable } from "svelte/store";
import en from "../../locales/en.json";
import pt from "../../locales/pt.json";
import schema from "../../locales/schema.json";
import { LocaleNames, SupportedLocales } from "../constants/locale";

/**
 * Identifiers of the locale schema.
 */
type LocaleTextID = keyof (typeof schema)["properties"];

/**
 * The standard locale in which the application is configured.
 */
const DEFAULT_LOCALE: SupportedLocales = SupportedLocales.EN;

/**
 * The standard locale name in which the application is configured.
 */
const DEFAULT_LOCALE_NAME: LocaleNames = LocaleNames.EN;

/**
 * Map of the supported locales with their respective configuration.
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
 * Retrieves a supported locale.
 * @param locale Source locale.
 * @returns `locale` if its a supported locale, otherwise {@link DEFAULT_LOCALE}.
 */
export function getSupportedLocale(locale: string): SupportedLocales {
	// Check if the selected locale in local storage is a supported locale.
	if (locale === SupportedLocales.EN || locale === SupportedLocales.PT) {
		return locale;
	}

	return DEFAULT_LOCALE;
}

/**
 * Retrieves the name of a supported locale.
 * @param locale Source locale.
 * @returns Name of `locale` if its a supported locale, otherwise {@link DEFAULT_LOCALE_NAME}.
 */
export function getSupportedLocaleName(locale: string): LocaleNames {
	switch (locale) {
		case "en":
			return LocaleNames.EN;
		case "pt":
			return LocaleNames.PT;
	}

	return DEFAULT_LOCALE_NAME;
}

/**
 * Inits a custom store for the application locale.
 * It synchronizes the store value with local storage.
 */
function _locale(): Writable<string> {
	const selectedLocale = localStorage.getItem("locale");

	let supportedLocale = DEFAULT_LOCALE;
	if (selectedLocale) {
		supportedLocale = getSupportedLocale(selectedLocale);
	}

	localStorage.setItem("locale", supportedLocale);

	// Set the supported locale in the lang attribute of the HTML element.
	document.documentElement.setAttribute("lang", supportedLocale);

	const store = writable(supportedLocale);
	const { subscribe, set } = store;

	return {
		subscribe,
		set(locale) {
			const supportedLocale = getSupportedLocale(locale);
			localStorage.setItem("locale", supportedLocale);

			// Set the supported locale in the lang attribute of the HTML element.
			document.documentElement.setAttribute("lang", supportedLocale);

			set(supportedLocale);
		},
		update(cb) {
			const updatedLocale = cb(get(store));

			const supportedLocale = getSupportedLocale(updatedLocale);
			localStorage.setItem("locale", supportedLocale);

			// Set the supported locale in the lang attribute of the HTML element.
			document.documentElement.setAttribute("lang", supportedLocale);

			set(supportedLocale);
		},
	};
}

/**
 * Locale store for reading and updating the locale of the application.
 */
export const locale = _locale();

/**
 * Derived store which translates a given text to the locale configured in the application.
 */
export const t = derived(locale, function ($locale) {
	/**
	 * Translates a given text to the locale configured in the application.
	 * @param id Identifier of the text to be localized.
	 * @param vars Map used to replace placeholder values inside a given text.
	 */
	return function (id: LocaleTextID, vars: Record<string, unknown> = {}) {
		return translate($locale, id, vars);
	};
});
