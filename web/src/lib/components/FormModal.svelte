<script lang="ts">
	import type { FormEventHandler } from "svelte/elements";
	import Modal from "./Modal.svelte";

	/**
	 * Indicates if modal is open.
	 */
	export let open: boolean;

	/**
	 * The modal title.
	 */
	export let title: string;

	/**
	 * Callback fired when the form is submitted.
	 */
	export let onSubmit: FormEventHandler<HTMLFormElement>;

	/**
	 * Callback fired when open state changes.
	 * @param open Modal open state.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Indicates whether the modal content is automatically configured with gutters.
	 * @default false
	 */
	export let gutters: boolean = false;

	/**
	 * The form element.
	 * @default null
	 */
	export let form: HTMLFormElement | null = null;

	/**
	 * Callback fired when backdrop is clicked.
	 * @default null
	 */
	export let onBackdropClick: (() => void) | null = null;
</script>

<form novalidate bind:this={form} on:submit|preventDefault={onSubmit}>
	<Modal {open} {gutters} {title} {onOpenChange} {onBackdropClick}>
		<slot />
		<svelte:fragment slot="actions">
			<slot name="actions" />
		</svelte:fragment>
	</Modal>
</form>
