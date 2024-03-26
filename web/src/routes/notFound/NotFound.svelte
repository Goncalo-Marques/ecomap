<script lang="ts">
	import { t } from "../../lib/utils/i8n";
</script>

<main>
	<h2>{$t("notFound.title")}</h2>
	<p>{$t("notFound.description")}</p>
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		justify-content: center;
		width: 25rem;
		margin: 0 auto;
		text-align: center;
	}
	h2 {
		font: var(--text-3xl-semibold);
	}
</style>
