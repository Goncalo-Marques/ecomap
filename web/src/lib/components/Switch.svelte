<script lang="ts">
	/**
	 * Value of toggle, true or false. Can be linked to a variable.
	 */
	export let checked: boolean = false;

	/**
	 * Function that will run when clicking on the component.
	 *
	 * @default None
	 */
	export let onClick = () => {};
</script>

<label class="switch">
	<input type="checkbox" bind:checked on:click={onClick} />
	<span class="slider round"></span>
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
		border-radius: 12px;
	}

	.slider.round:before {
		border-radius: 50%;
	}
</style>
