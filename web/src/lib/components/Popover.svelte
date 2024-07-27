<script lang="ts">
	/**
	 * The popover element.
	 * @default null
	 */
	export let popover: HTMLElement | null = null;

	/**
	 * The identifier of the popover element.
	 * @default "popover"
	 */
	export let id: string = "popover";

	/**
	 * The position of the popover element relative to the trigger element.
	 * @default "left"
	 */
	export let align: "left" | "middle" | "right" = "left";

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

		// Set the popover display to block to get its dimensions.
		popover.style.display = "block";

		const viewportTriggerOffset = triggerElement.getBoundingClientRect();

		switch (align) {
			case "left":
				popover.style.left = `${viewportTriggerOffset.left}px`;
				popover.style.top = `${viewportTriggerOffset.bottom}px`;
				break;

			case "middle":
				popover.style.left = `${viewportTriggerOffset.left + (viewportTriggerOffset.width - popover.offsetWidth) / 2}px`;
				popover.style.top = `${viewportTriggerOffset.bottom}px`;
				break;

			case "right":
				popover.style.left = `${viewportTriggerOffset.left + viewportTriggerOffset.width - popover.offsetWidth}px`;
				popover.style.top = `${viewportTriggerOffset.bottom}px`;
				break;
		}

		// Clear the display property after the popover is open.
		popover.style.display = "";
	}
</script>

<button
	class="flex items-center justify-center"
	popovertarget={id}
	on:click={handleTriggerClick}
>
	<slot name="trigger" />
</button>

<div
	bind:this={popover}
	role="menu"
	class="m-0 rounded bg-white p-0 shadow-md"
	popover="auto"
	{id}
>
	<slot />
</div>
