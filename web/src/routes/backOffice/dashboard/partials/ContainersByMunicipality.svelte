<script lang="ts">
	import { Chart } from "chart.js";
	import Card from "../../components/Card.svelte";
	import type { Container } from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { getCssVariable } from "../../../../lib/utils/cssVars";
	import { t } from "../../../../lib/utils/i8n";

	export let containersPromise: Promise<Container[]>;

	let loading = false;

	function makeChart(containers: Container[]) {
		const canvasElement = document.getElementById(
			"container-municipality-chart",
		);
		if (!(canvasElement instanceof HTMLCanvasElement)) {
			return;
		}

		const data = new Map<string, number>();
		for (const container of containers) {
			const name =
				container.geoJson.properties.municipalityName ??
				$t("location.unknownWay");
			const containerMunicipalityAmount = data.get(name) ?? 0;
			data.set(name, containerMunicipalityAmount + 1);
		}

		new Chart(canvasElement, {
			type: "bar",
			data: {
				labels: Array.from(data.keys()),
				datasets: [
					{
						data: Array.from(data.values()),
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

	containersPromise.then(containers => makeChart(containers));
</script>

<Card element="article" class="containers-municipality-card">
	<h2>Contentores por município</h2>
	<canvas
		id="container-municipality-chart"
		style={loading ? "display: none;" : ""}
	/>
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>Erro</p>
	{/await}
</Card>

<style>
	:global(.containers-municipality-card) {
		grid-area: containersByMunicipality;
	}

	h2 {
		font: var(--text-xl-semibold);
	}
</style>
