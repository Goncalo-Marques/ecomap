<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";

	/**
	 * Indicates if switch is checked.
	 */
	export let checked: boolean = false;

	/**
	 * Callback fired when clicking the component.
	 *
	 * @default null
	 */
	export let onClick: MouseEventHandler<HTMLInputElement> | null = null;
</script>

<label class="switch">
	<input type="checkbox" bind:checked on:click={onClick} />
	<span class="slider round" />
</label>

<style>
	.switch {
		margin-left: auto;

		position: relative;
		display: inline-block;
		width: 2.5rem;
		height: 1.25rem;
	}

	.switch input {
		opacity: 0;
	}

	.slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		transition: 0.2s;
	}

	.slider:before {
		position: absolute;
		content: "";
		height: 0.75rem;
		width: 0.75rem;
		bottom: 0.25rem;
		transition: 0.2s;
	}

	input:checked + .slider {
		outline: 0;
		background-color: var(--green-100);
	}

	input:focus + .slider {
		box-shadow: 0 0 1px var(--green-100);
	}

	input:checked + .slider:before {
		transform: translateX(1.375rem);
		background-color: var(--green-800);
	}

	input + .slider::before {
		transform: translateX(0.25rem);
		background-color: var(--gray-400);
	}

	input + .slider {
		outline: 1px solid var(--gray-400);
	}

	.slider.round {
		border-radius: 0.75rem;
	}

	.slider.round:before {
		border-radius: 50%;
	}
</style>
