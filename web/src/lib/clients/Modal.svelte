<script lang="ts">
	import { onMount } from "svelte";
	import IconButton from "../components/IconButton.svelte";

	export let open: boolean;

	export let title: string;

	export let gutters: boolean = false;

	export let onClose: () => void;

	let dialog: HTMLDialogElement;

	function closeModal() {
		onClose();
		dialog.close();
	}

	function handleClick(e: Event) {
		if (e.target === dialog) {
			closeModal();
		}
	}

	onMount(() => {
		dialog.addEventListener("click", handleClick);
	});

	$: if (!dialog?.open && open) {
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
	<footer></footer>
</dialog>

<style>
	dialog {
		padding: 0;
		margin: auto;
		min-width: 32rem;
		border: unset;
		border-radius: 0.25rem;
		box-shadow: var(--shadow-md);
	}

	dialog::backdrop {
		background-color: black;
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
		padding: 1.5rem;
	}
</style>
