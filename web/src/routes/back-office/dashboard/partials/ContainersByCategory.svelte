<script lang="ts">
	import { onDestroy } from "svelte";

	import { Chart } from "chart.js";

	import Spinner from "$lib/components/Spinner.svelte";
	import { getCssVariable } from "$lib/utils/cssVars";
	import { t } from "$lib/utils/i8n";

	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Card from "../../components/Card.svelte";

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
	let canvas: HTMLCanvasElement | undefined;

	/**
	 * The chart being rendered.
	 */
	let chart: Chart<"doughnut"> | undefined;

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
		// Exit if canvas is not bound to DOM element.
		if (!canvas) {
			return;
		}

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

		chart = new Chart(canvas, {
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

	onDestroy(() => {
		// Destroy chart before the component destruction.
		chart?.destroy();
	});

	containersPromise.then(containers => buildChart(containers));
</script>

<Card element="article" class="[grid-area:containersByCategory]">
	<h2 class="text-xl font-semibold">{$t("dashboard.containersByCategory")}</h2>
	<canvas bind:this={canvas} style:display={loading ? "none" : ""} />
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>{$t("error.unexpected.title")}</p>
	{/await}
</Card>
