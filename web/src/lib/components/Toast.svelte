<script lang="ts">
	import { onMount, setContext } from "svelte";
	import type { ToastContext, ToastOptions } from "../../domain/toast";
	import Icon from "./Icon.svelte";
	import { writable } from "svelte/store";
	import { setToastContext } from "../contexts/toast";

	let toast: HTMLElement;

	let toastOptions = writable<ToastOptions>({
		title: "Toast",
		type: "success",
	});

	/**
	 * The icons to be shown for each toast type.
	 */
	const icons: Record<ToastOptions["type"], string> = {
		success: "check",
		error: "error",
	};

	/**
	 * Duration of the toast in milliseconds.
	 */
	const TOAST_DURATION = 5000;

	/**
	 * Shows the toast element.
	 */
	function showToast() {
		// Appends toast to end of body element in the DOM.
		document.body.appendChild(toast);
		toast.showPopover();

		// Removes the toast again after the given toast duration.
		setTimeout(() => {
			toast.hidePopover();
			toast.remove();
		}, TOAST_DURATION);
	}

	// Sets the toast context to enable interaction in the app components.
	setToastContext({
		show(options: ToastOptions) {
			toastOptions.set(options);
			showToast();
		},
	});

	onMount(() => {
		// Removes the toast from the current position in the DOM.
		// The toast is only in the DOM whenever the show() method is called.
		toast.remove();
	});
</script>

<slot />
<div
	class="toast"
	popover="manual"
	data-type={$toastOptions.type}
	bind:this={toast}
>
	<Icon name={icons[$toastOptions.type]} size="small" />
	<div class="content">
		<h3>{$toastOptions.title}</h3>
		{#if $toastOptions.description}
			<p>{$toastOptions.description}</p>
		{/if}
	</div>
</div>

<style>
	.toast {
		display: flex;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		border-radius: 0.25rem;
		background-color: var(--white);
		box-shadow: var(--shadow-lg);
		border-style: solid;
		border-width: 1px;

		&[data-type="success"] {
			background-color: var(--green-50);
			color: var(--green-700);
			border-color: var(--green-300);

			& p {
				color: var(--gray-900);
			}
		}
		&[data-type="error"] {
			background-color: var(--red-50);
			color: var(--red-700);
			border-color: var(--red-300);

			& p {
				color: var(--gray-900);
			}
		}

		& .content {
			display: flex;
			flex-direction: column;
		}

		& h3 {
			font: var(--text-base-semibold);
		}

		& p {
			font: var(--text-sm-regular);
		}
	}

	:popover-open {
		position: absolute;
		inset: unset;
		left: 1rem;
		bottom: 1rem;
	}
</style>
