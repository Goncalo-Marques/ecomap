<script lang="ts">
	import Icon from "$lib/components/Icon.svelte";
	import Modal from "$lib/components/Modal.svelte";
	import { t } from "$lib/utils/i8n";

	const team = [
		{
			name: "Gonçalo Dias",
			email: "goncalo.gjad@gmail.com",
			linkedInURL: "https://www.linkedin.com/in/gjdias",
		},
		{
			name: "Gonçalo Marques",
			email: "goncalomiguelsilvamarques@gmail.com",
			linkedInURL: "https://www.linkedin.com/in/goncalosilvamarques",
		},
		{
			name: "João Pinheiro",
			email: "joaotomasp@gmail.com",
			linkedInURL: "https://www.linkedin.com/in/joaotomaspinheiro",
		},
	];

	/**
	 * Indicates if the support modal is open.
	 */
	let openSupportModal: boolean;
</script>

<button
	class="flex w-full justify-items-center gap-2 rounded px-3 py-2 transition-colors hover:bg-green-50 hover:text-green-800 active:bg-green-100 active:text-green-800 data-[selected=true]:bg-green-100 data-[selected=true]:font-semibold data-[selected=true]:text-green-800"
	on:click={() => (openSupportModal = true)}
>
	<Icon name="headset_mic" />
	<span>{$t("support")}</span>
</button>

<Modal
	open={openSupportModal}
	gutters
	title={$t("support")}
	onOpenChange={open => (openSupportModal = open)}
>
	<ul class="flex flex-col gap-2">
		{#each team as member}
			<li class="flex flex-1 items-center justify-between gap-2">
				<div class="flex flex-col">
					<span>{member.name}</span>
					<span class="text-sm text-gray-600">{member.email}</span>
				</div>
				<a
					class="flex size-8 items-center justify-center rounded-full bg-[#007ebb]"
					target="_blank"
					href={member.linkedInURL}
				>
					<img class="size-5" src="/images/linkedin.svg" alt="LinkedIn logo" />
				</a>
			</li>
		{/each}
	</ul>
</Modal>
