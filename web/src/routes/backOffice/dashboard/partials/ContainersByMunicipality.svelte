<script lang="ts">
	import { onDestroy } from "svelte";
	import { Chart } from "chart.js";
	import Card from "../../components/Card.svelte";
	import type { Container } from "../../../../domain/container";
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
	let canvas: HTMLCanvasElement | undefined;

	/**
	 * The chart being rendered.
	 */
	let chart: Chart<"bar"> | undefined;

	/**
	 * Builds a chart with the containers.
	 * @param containers Containers.
	 */
	function buildChart(containers: Container[]) {
		// Exit if canvas is not bound to DOM element.
		if (!canvas) {
			return;
		}

		// Build a map with the amount of containers per municipality.
		const containersPerMunicipality = new Map<string, number>();
		for (const container of containers) {
			const municipalityName = container.geoJson.properties.municipalityName;
			if (!municipalityName) {
				continue;
			}

			const amount = containersPerMunicipality.get(municipalityName) ?? 0;
			containersPerMunicipality.set(municipalityName, amount + 1);
		}

		// Get chart labels.
		const labels = Array.from(containersPerMunicipality.keys());

		// Get chart data.
		const data = Array.from(containersPerMunicipality.values());

		chart = new Chart(canvas, {
			type: "bar",
			data: {
				labels,
				datasets: [
					{
						data,
						backgroundColor: getCssVariable("--green-700"),
					},
				],
			},
			options: {
				plugins: {
					legend: {
						display: false,
					},
				},
				scales: {
					y: {
						beginAtZero: true,
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

<Card element="article" class="[grid-area:containersByMunicipality]">
	<h2 class="text-xl font-semibold">
		{$t("dashboard.containersByMunicipality")}
	</h2>
	<canvas bind:this={canvas} style:display={loading ? "none" : ""} />
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>{$t("error.unexpected.title")}</p>
	{/await}
</Card>
