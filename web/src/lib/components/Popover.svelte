<script lang="ts">
	/**
	 * The popover element.
	 * @default null
	 */
	export let popover: HTMLElement | null = null;

	/**
	 * Handles click event on trigger element of the popover.
	 * @param e Click event.
	 */
	function handleTriggerClick(e: Event) {
		const triggerElement = e.currentTarget as HTMLButtonElement;

		const popover = triggerElement.popoverTargetElement as HTMLElement | null;
		if (!popover) {
			return;
		}

		const viewportOffset = triggerElement.getBoundingClientRect();
		popover.style.left = `${viewportOffset.left}px`;
		popover.style.top = `${viewportOffset.bottom}px`;
	}
</script>

<button popovertarget="popover" on:click={handleTriggerClick}>
	<slot name="trigger" />
</button>

<div
	bind:this={popover}
	role="menu"
	class="dropdown"
	popover="auto"
	id="popover"
>
	<slot />
</div>

<style>
	button {
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.dropdown {
		background-color: var(--white);
		box-shadow: var(--shadow-md);
		padding: 1rem;
		border-radius: 0.25rem;
	}
</style>
