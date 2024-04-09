<script lang="ts">

	/**
	 * Value of toggle, true or false. Can be linked to a variable.
	 */
    export let checked: boolean
    
	/**
	 * Function that will run when clicking on the component.
	 * 
	 * @default None
	 */
	export let onClick = () => {}
</script>

<label class="switch">
    <input type="checkbox" bind:checked={checked} on:click={onClick}/>
    <span class="slider round"></span>
</label>

<style>

.switch { 
		box-sizing: border-box;
		margin-left: auto;
		
		position: relative;
		display: inline-block;
		width: 40px;
		height: 20px;
	}

	.switch input {
		opacity: 0;
		width: 0;
		height: 0;
	}

	.slider {
		position: absolute;
		cursor: pointer;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: var(--white);
		-webkit-transition: 0.2s;
		transition: 0.2s;
	}

	.slider:before {
		position: absolute;
		content: "";
		height: 12px;
		width: 12px;
		bottom: 4px;

		-webkit-transition: 0.2s;
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
		-webkit-transform: translateX(22px);
		-ms-transform: translateX(22px);
		transform: translateX(22px);
		background-color: var(--green-800);

	}

	input + .slider::before{
		transform: translateX(4px);
		background-color: var(--gray-400);
	}

	input + .slider {
		outline: 1px solid var(--gray-400);
	}

	/* Rounded sliders */
	.slider.round {
		border-radius: 12px;
	}

	.slider.round:before {
		border-radius: 50%;
	}
</style>
