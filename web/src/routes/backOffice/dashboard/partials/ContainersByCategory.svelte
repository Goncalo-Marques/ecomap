<script lang="ts">
	import { Chart } from "chart.js";
	import Card from "../../components/Card.svelte";
	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { getCssVariable } from "../../../../lib/utils/cssVars";
	import { t } from "../../../../lib/utils/i8n";

	/**
	 * The promise with the containers.
	 */
	export let containersPromise: Promise<Container[]>;

	/**
	 * Indicates whether the containers are being loaded.
	 */
	let loading = true;

	/**
	 * The canvas element where the chart is rendered.
	 */
	let canvas: HTMLCanvasElement;

	/**
	 * The color for each respective container category.
	 */
	const categoryColors: Record<ContainerCategory, string> = {
		general: getCssVariable("--sky-300"),
		glass: getCssVariable("--lime-600"),
		hazardous: getCssVariable("--amber-500"),
		metal: getCssVariable("--indigo-400"),
		organic: getCssVariable("--orange-800"),
		paper: getCssVariable("--blue-400"),
		plastic: getCssVariable("--yellow-400"),
	};

	/**
	 * Builds a chart with the containers.
	 * @param containers Containers.
	 */
	function buildChart(containers: Container[]) {
		// Build a map with the amount of containers per container category.
		const containerAmountPerCategory = new Map<ContainerCategory, number>();
		for (const container of containers) {
			const amount = containerAmountPerCategory.get(container.category) ?? 0;
			containerAmountPerCategory.set(container.category, amount + 1);
		}

		// Get container categories.
		const containerCategories = Array.from(containerAmountPerCategory.keys());

		// Get chart labels.
		const labels = containerCategories.map(category =>
			$t(`containers.category.${category}`),
		);

		// Get chart data.
		const data = Array.from(containerAmountPerCategory.values());

		// Get background color for each container category.
		const backgroundColor = containerCategories.map(
			category => categoryColors[category],
		);

		new Chart(canvas, {
			type: "doughnut",
			data: {
				labels,
				datasets: [
					{
						data,
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

	containersPromise.then(containers => buildChart(containers));
</script>

<Card element="article" class="containers-category-card">
	<h2>{$t("dashboard.containersByCategory")}</h2>
	<canvas bind:this={canvas} style:display={loading ? "none" : ""} />
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>{$t("error.unexpected.title")}</p>
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
