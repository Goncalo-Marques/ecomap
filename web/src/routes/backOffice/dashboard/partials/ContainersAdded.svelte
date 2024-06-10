<script lang="ts">
	import { Chart } from "chart.js";
	import Card from "../../components/Card.svelte";
	import type { Container } from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { getCssVariable } from "../../../../lib/utils/cssVars";

	export let containersPromise: Promise<Container[]>;

	let loading = false;

	function getMonths(): Map<number, string> {
		const today = new Date();
		const monthsNamesMap = new Map<number, string>();

		for (let i = 0; i < 12; i++) {
			const date = new Date(today.getFullYear(), i, 1);
			const month = date.toLocaleString("default", { month: "short" });
			monthsNamesMap.set(i, month);
		}

		return monthsNamesMap;
	}

	function makeChart(containers: Container[]) {
		const canvasElement = document.getElementById("container-added-chart");
		if (!(canvasElement instanceof HTMLCanvasElement)) {
			return;
		}

		const containersPerMonth = new Map<number, number>();
		const monthNames = getMonths();
		for (const container of containers) {
			const createdAt = new Date(container.createdAt);
			const monthDateInMilliseconds = new Date(
				createdAt.getFullYear(),
				createdAt.getMonth(),
				1,
			).valueOf();

			const amount = containersPerMonth.get(monthDateInMilliseconds) ?? 0;
			containersPerMonth.set(monthDateInMilliseconds, amount + 1);
		}

		const labels = Array.from(containersPerMonth.keys()).map(dateMs => {
			const date = new Date(dateMs);
			const monthName = monthNames.get(date.getMonth())!;
			const year = date.getFullYear();

			return `${monthName} ${year}`;
		});

		let data = Array.from(containersPerMonth.values());
		data = data.reduce((acc, amount, idx) => {
			if (idx > 0) {
				acc[idx] = acc[idx - 1] + amount;
			}

			return acc;
		}, data);

		new Chart(canvasElement, {
			type: "line",
			data: {
				labels,
				datasets: [
					{
						data,
						borderColor: getCssVariable("--green-700"),
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

<Card element="article" class="containers-added-card">
	<h2>Contentores adicionados</h2>
	<canvas id="container-added-chart" style={loading ? "display: none;" : ""} />
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>Erro</p>
	{/await}
</Card>

<style>
	:global(.containers-added-card) {
		grid-area: containersAdded;
	}

	h2 {
		font: var(--text-xl-semibold);
	}
</style>
