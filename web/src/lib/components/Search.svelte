<script lang="ts">
	import type {
		ChangeEventHandler,
		HTMLInputAttributes,
	} from "svelte/elements";
	import Icon from "./Icon.svelte";

	/**
	 * The hint for form autofill feature.
	 * @default "off"
	 */
	export let autocomplete: HTMLInputAttributes["autocomplete"] = "off";

	/**
	 * Callback fired when input value changes.
	 * @default null
	 */
	export let onInput: ChangeEventHandler<HTMLInputElement> | null = null;

	/**
	 * The text that appears in the form control when it has no value set.
	 * @default null
	 */
	export let placeholder: HTMLInputAttributes["placeholder"] = null;

	/**
	 * The value of the search input.
	 * @default null
	 */
	export let value: string | null = null;
</script>

<label>
	<div class="search-icon">
		<Icon name="search" size="small" />
	</div>
	<input
		{autocomplete}
		{placeholder}
		type="search"
		{value}
		on:input={onInput}
	/>
</label>

<style>
	label {
		position: relative;
		display: block;

		&:focus-within {
			& input {
				padding-left: 0.5rem;
			}

			& .search-icon {
				display: none;
			}
		}
	}

	input {
		position: relative;
		width: 100%;
		padding: 0.375rem 0.5rem 0.375rem 2rem;
		border: 1px solid var(--gray-300);
		border-radius: 0.25rem;
		color: var(--gray-900);

		&::placeholder {
			color: var(--gray-400);
		}

		&::-webkit-search-cancel-button {
			appearance: none;
			cursor: pointer;
			position: relative;
			margin: 0;
			height: 1.25rem;
			width: 1.25rem;
			background-image: url("data:image/svg+xml,%3Csvg%20width%3D%2224%22%20height%3D%2224%22%20viewBox%3D%220%200%2024%2024%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Cpath%20d%3D%22M12%2013.4L7.09999%2018.3C6.91665%2018.4834%206.68332%2018.575%206.39999%2018.575C6.11665%2018.575%205.88332%2018.4834%205.69999%2018.3C5.51665%2018.1167%205.42499%2017.8834%205.42499%2017.6C5.42499%2017.3167%205.51665%2017.0834%205.69999%2016.9L10.6%2012L5.69999%207.10005C5.51665%206.91672%205.42499%206.68338%205.42499%206.40005C5.42499%206.11672%205.51665%205.88338%205.69999%205.70005C5.88332%205.51672%206.11665%205.42505%206.39999%205.42505C6.68332%205.42505%206.91665%205.51672%207.09999%205.70005L12%2010.6L16.9%205.70005C17.0833%205.51672%2017.3167%205.42505%2017.6%205.42505C17.8833%205.42505%2018.1167%205.51672%2018.3%205.70005C18.4833%205.88338%2018.575%206.11672%2018.575%206.40005C18.575%206.68338%2018.4833%206.91672%2018.3%207.10005L13.4%2012L18.3%2016.9C18.4833%2017.0834%2018.575%2017.3167%2018.575%2017.6C18.575%2017.8834%2018.4833%2018.1167%2018.3%2018.3C18.1167%2018.4834%2017.8833%2018.575%2017.6%2018.575C17.3167%2018.575%2017.0833%2018.4834%2016.9%2018.3L12%2013.4Z%22%20fill%3D%22%23111827%22%2F%3E%3C%2Fsvg%3E");
			background-size: contain;
		}
	}

	.search-icon {
		pointer-events: none;
		position: absolute;
		left: 0.5rem;
		top: 50%;
		transform: translateY(-50%);
		display: flex;
		justify-content: center;
		align-items: center;
		color: var(--gray-400);
	}
</style>
