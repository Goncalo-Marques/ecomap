<script lang="ts">
	import { Link } from "svelte-routing";
	import Button from "./Button.svelte";

	/**
	 * The title of the bottom sheet.
	 */
	export let title: string;

	/**
	 * The page link to the resource details.
	 */
	export let resourceLink: string | null = null;
</script>

<div class="absolute bottom-0 left-0 flex w-full px-10">
	<div class="flex max-h-80 flex-1 flex-col gap-4 rounded-t bg-white shadow-lg">
		<div class="flex px-8 pt-8">
			<h2 class="flex-1 text-2xl font-semibold">{title}</h2>
			{#if resourceLink}
				<Link to={resourceLink} class="contents">
					<Button variant="tertiary" startIcon="open_in_new" />
				</Link>
			{/if}
		</div>
		<div
			class="grid grid-cols-[repeat(auto-fill,minmax(max(13rem,(100%-0.5rem)/4),1fr))] gap-2 overflow-auto px-8 pb-8"
		>
			<slot />
		</div>
	</div>
</div>
