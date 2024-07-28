<script lang="ts">
	import { onMount } from "svelte";

	import IconButton from "./IconButton.svelte";

	/**
	 * Indicates if modal is open.
	 */
	export let open: boolean;

	/**
	 * The modal title.
	 */
	export let title: string;

	/**
	 * Indicates whether the modal content is automatically configured with gutters.
	 * @default false
	 */
	export let gutters: boolean = false;

	/**
	 * Callback fired when open state changes.
	 * @param open Modal open state.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Callback fired when backdrop is clicked.
	 * @default null
	 */
	export let onBackdropClick: (() => void) | null = null;

	/**
	 * Dialog element.
	 */
	let dialog: HTMLDialogElement;

	/**
	 * Indicates whether the modal contains actions.
	 * If the modal contains actions, a footer is added to the modal.
	 */
	let containsActions: boolean = !!$$props.$$slots.actions;

	/**
	 * Closes the modal.
	 */
	function closeModal() {
		onOpenChange(false);
		dialog.close();
	}

	/**
	 * Handles the click event of the dialog element.
	 * @param e Click event.
	 */
	function handleDialogClick(e: Event) {
		// Check if the element that was pressed is the dialog. If it is, it means that the click was performed outside the modal.
		if (e.target === dialog) {
			closeModal();
			onBackdropClick?.();
		}
	}

	onMount(() => {
		dialog.addEventListener("click", handleDialogClick);
	});

	// Watch for changes in the modal open state to show or hide the modal.
	$: if (!dialog?.open && open) {
		onOpenChange(true);
		dialog.showModal();
	} else if (dialog?.open && !open) {
		closeModal();
	}
</script>

<dialog
	class="m-auto max-h-[90vh] min-w-[32rem] max-w-[90vw] overflow-hidden rounded shadow-md backdrop:bg-black backdrop:opacity-50"
	bind:this={dialog}
>
	<header class="border-b border-gray-300 py-6">
		<div class="flex items-center gap-2 px-6">
			<h2 class="flex-1 text-2xl font-semibold">{title}</h2>
			<IconButton icon="close" onClick={closeModal} />
		</div>
	</header>
	<section style:padding={gutters ? "1.5rem" : ""}>
		<slot />
	</section>
	{#if containsActions}
		<footer class="flex items-center justify-end gap-2 p-6">
			<slot name="actions" />
		</footer>
	{/if}
</dialog>
