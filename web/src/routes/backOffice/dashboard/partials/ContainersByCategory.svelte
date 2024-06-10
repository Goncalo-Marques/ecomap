<script lang="ts">
	import { Chart } from "chart.js";
	import Card from "../../components/Card.svelte";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { getCssVariable } from "../../../../lib/utils/cssVars";
	import { t } from "../../../../lib/utils/i8n";

	export let containersPromise: Promise<Container[]>;

	let loading = false;

	const categoryColors: Record<ContainerCategory, string> = {
		general: getCssVariable("--sky-300"),
		glass: getCssVariable("--lime-600"),
		hazardous: getCssVariable("--amber-500"),
		metal: getCssVariable("--indigo-400"),
		organic: getCssVariable("--orange-800"),
		paper: getCssVariable("--blue-400"),
		plastic: getCssVariable("--yellow-400"),
	};

	function makeChart(containers: Container[]) {
		const canvasElement = document.getElementById("container-category-chart");
		if (!(canvasElement instanceof HTMLCanvasElement)) {
			return;
		}

		const data = new Map<ContainerCategory, number>();
		for (const container of containers) {
			const containerCategoryAmount = data.get(container.category) ?? 0;
			data.set(container.category, containerCategoryAmount + 1);
		}

		const backgroundColor = Array.from(data.keys()).map(
			category => categoryColors[category],
		);

		new Chart(canvasElement, {
			type: "doughnut",
			data: {
				labels: Array.from(data.keys()).map(category =>
					$t(`containers.category.${category}`),
				),
				datasets: [
					{
						data: Array.from(data.values()),
						backgroundColor,
					},
				],
			},
			options: {
				plugins: {
					legend: {
						position: "right",
					},
				},
			},
		});

		loading = false;
	}

	containersPromise.then(containers => makeChart(containers));
</script>

<Card element="article" class="containers-category-card">
	<h2>Contentores por categoria</h2>
	<canvas
		id="container-category-chart"
		style={loading ? "display: none;" : ""}
	/>
	{#await containersPromise}
		<Spinner />
	{:catch error}
		<p>{error}</p>
	{/await}
</Card>

<style>
	:global(.containers-category-card) {
		grid-area: containersByCategory;
	}

	h2 {
		font: var(--text-xl-semibold);
	}
</style>
