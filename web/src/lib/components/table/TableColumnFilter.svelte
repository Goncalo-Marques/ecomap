<script lang="ts">
	import Icon from "../Icon.svelte";

	function handleFilterClick(e: Event) {
		const filterButton = e.currentTarget as HTMLButtonElement;
		const filterDropdown = filterButton.parentElement;

		if (!filterDropdown) {
			return;
		}

		const { state } = filterDropdown.dataset;

		switch (state) {
			case "closed":
				filterDropdown.dataset.state = "open";
				break;
			case "open":
				filterDropdown.dataset.state = "closed";
				break;
		}
	}
</script>

<div class="dropdown" data-state="closed">
	<button on:click={handleFilterClick}>
		<Icon name="filter_alt" size="x-small" />
	</button>
	<div class="dropdown-content">
		<li>Option 1</li>
		<li>Option 2</li>
		<li>Option 3</li>
	</div>
</div>

<style>
	button {
		display: contents;
	}

	.dropdown {
		position: relative;
		display: inline-block;

		&[data-state="closed"] {
			& .dropdown-content {
				display: none;
			}
		}
	}

	.dropdown-content {
		display: flex;
		flex-direction: column;
		position: absolute;
		background-color: var(--white);
		box-shadow: var(--shadow-md);
		padding: 1rem;
		min-width: 8rem;
		border-radius: 0.25rem;
		z-index: 9999;
	}
</style>
