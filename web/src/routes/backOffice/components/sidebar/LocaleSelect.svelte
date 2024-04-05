<script lang="ts">
	import { getSupportedLocale, locale } from "../../../../lib/utils/i8n";

	/**
	 * Handles on change event of the locale select input.
	 * @param e Select on change event.
	 */
	function handleLocaleChange(e: Event) {
		const select = e.currentTarget as HTMLSelectElement;
		locale.set(select.value);
	}
</script>

<select value={$locale} on:change={handleLocaleChange}>
	<option value="en">English</option>
	<option value="pt">Português</option>
</select>

<style>
	select {
		border: 1px solid var(--gray-300);
		padding: 0.5rem 0.75rem;
		border-radius: 0.25rem;
	}
</style>
