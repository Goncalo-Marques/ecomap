<script lang="ts">
	import { onMount } from "svelte";
	import type { ToastOptions } from "../../domain/toast";
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
	class="flex max-w-[36rem] gap-2 rounded border bg-white px-4 py-3 shadow-lg data-[type=error]:border-red-300 data-[type=success]:border-green-300 data-[type=error]:bg-red-50 data-[type=success]:bg-green-50 data-[type=error]:text-red-700 data-[type=success]:text-green-700 [&:popover-open]:absolute [&:popover-open]:inset-[unset] [&:popover-open]:bottom-4 [&:popover-open]:left-4"
	popover="manual"
	data-type={$toastOptions.type}
	bind:this={toast}
>
	<Icon name={icons[$toastOptions.type]} size="small" />
	<div class="flex flex-col">
		<h3 class="font-semibold">{$toastOptions.title}</h3>
		{#if $toastOptions.description}
			<p class="text-sm text-gray-900">
				{$toastOptions.description}
			</p>
		{/if}
	</div>
</div>
