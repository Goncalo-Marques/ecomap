<script lang="ts">
	import Icon from "../../../../lib/components/Icon.svelte";
	import Modal from "../../../../lib/components/Modal.svelte";
	import { t } from "../../../../lib/utils/i8n";

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

<button on:click={() => (openSupportModal = true)}>
	<Icon name="headset_mic" />
	<span>{$t("support")}</span>
</button>

<Modal
	open={openSupportModal}
	gutters
	title={$t("support")}
	onOpenChange={open => (openSupportModal = open)}
>
	<ul>
		{#each team as member}
			<li>
				<div class="member-details">
					<span class="member-details-name">{member.name}</span>
					<span class="member-details-email">{member.email}</span>
				</div>
				<a target="_blank" href={member.linkedInURL}>
					<img src="/images/linkedin.svg" alt="LinkedIn logo" />
				</a>
			</li>
		{/each}
	</ul>
</Modal>

<style>
	ul {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		& li {
			flex: 1;
			display: flex;
			justify-content: space-between;
			align-items: center;
			gap: 0.5rem;

			& .member-details {
				display: flex;
				flex-direction: column;

				& .member-details-name {
					font: var(--text-base-regular);
				}

				& .member-details-email {
					font: var(--text-sm-regular);
					color: var(--gray-600);
				}
			}

			& a {
				width: 2rem;
				height: 2rem;
				display: flex;
				justify-content: center;
				align-items: center;
				border-radius: 50%;
				background-color: #007ebb;

				& img {
					width: 1.25rem;
					height: 1.25rem;
				}
			}
		}
	}
</style>
