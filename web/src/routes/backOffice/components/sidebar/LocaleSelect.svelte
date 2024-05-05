<script lang="ts">
	import Option from "../../../../lib/components/Option.svelte";
	import Select from "../../../../lib/components/Select.svelte";
	import { locale } from "../../../../lib/utils/i8n";

	/**
	 * Handles on change event of the locale select input.
	 * @param e Select on change event.
	 */
	function handleLocaleChange(e: Event) {
		const select = e.currentTarget as HTMLSelectElement;
		locale.set(select.value);
	}
</script>

<Select value={$locale} onChange={handleLocaleChange}>
	<Option value="en">English</Option>
	<Option value="pt">Português</Option>
</Select>
