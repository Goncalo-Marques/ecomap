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

<dialog bind:this={dialog}>
	<header>
		<div class="header-content">
			<h2>{title}</h2>
			<IconButton icon="close" onClick={closeModal} />
		</div>
	</header>
	<section style:padding={gutters ? "1.5rem" : ""}>
		<slot />
	</section>
	{#if containsActions}
		<footer>
			<slot name="actions" />
		</footer>
	{/if}
</dialog>

<style>
	dialog {
		padding: 0;
		margin: auto;
		min-width: 32rem;
		max-width: 90vw;
		max-height: 90vh;
		border: unset;
		border-radius: 0.25rem;
		box-shadow: var(--shadow-md);
		overflow: hidden;
	}

	dialog::backdrop {
		background-color: var(--black);
		opacity: 0.5;
	}

	header {
		padding-block: 1.5rem;
		border-bottom: 1px solid var(--gray-300);

		& > .header-content {
			display: flex;
			align-items: center;
			gap: 0.5rem;
			padding-inline: 1.5rem;

			& > h2 {
				font: var(--text-2xl-semibold);
				flex: 1;
			}
		}
	}

	footer {
		display: flex;
		justify-content: flex-end;
		align-items: center;
		gap: 0.5rem;
		padding: 1.5rem;
	}
</style>
