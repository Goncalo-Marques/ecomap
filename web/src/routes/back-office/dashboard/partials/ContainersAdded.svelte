<script lang="ts">
	import { onDestroy } from "svelte";
	import { get } from "svelte/store";

	import { Chart } from "chart.js";

	import Spinner from "$lib/components/Spinner.svelte";
	import { getColorWithOpacity } from "$lib/utils/color";
	import { getCssVariable } from "$lib/utils/cssVars";
	import { locale, t } from "$lib/utils/i8n";

	import type { Container } from "../../../../domain/container";
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
	let chart: Chart<"line"> | undefined;

	/**
	 * Retrieves a map with the month index and the corresponding names.
	 * @returns Map with the month index and the corresponding name.
	 */
	function getMonths(): Map<number, string> {
		const today = new Date();
		const monthsNamesMap = new Map<number, string>();

		for (let i = 0; i < 12; i++) {
			const date = new Date(today.getFullYear(), i, 1);
			const month = date.toLocaleString(get(locale), { month: "short" });
			monthsNamesMap.set(i, month);
		}

		return monthsNamesMap;
	}

	/**
	 * Builds a chart with the containers.
	 * @param containers Containers.
	 */
	function buildChart(containers: Container[]) {
		// Exit if canvas is not bound to DOM element.
		if (!canvas) {
			return;
		}

		// Build a map with the amount of containers added per month.
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

		// Get chart labels.
		const labels = Array.from(containersPerMonth.keys()).map(dateMs => {
			const date = new Date(dateMs);
			const monthName = monthNames.get(date.getMonth())!;
			const year = date.getFullYear();

			return `${monthName} ${year}`;
		});

		// Get chart data.
		let data = Array.from(containersPerMonth.values());
		data = data.reduce((acc, amount, idx) => {
			if (idx > 0) {
				acc[idx] = acc[idx - 1] + amount;
			}

			return acc;
		}, data);

		chart = new Chart(canvas, {
			type: "line",
			data: {
				labels,
				datasets: [
					{
						data,
						borderColor: getCssVariable("--green-700"),
						backgroundColor: getColorWithOpacity(
							getCssVariable("--green-700"),
							0.2,
						),
						pointBackgroundColor: getCssVariable("--white"),
						fill: "origin",
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

<Card element="article" class="[grid-area:containersAdded]">
	<h2 class="text-xl font-semibold">{$t("dashboard.containersAdded")}</h2>
	<canvas bind:this={canvas} style:display={loading ? "none" : ""} />
	{#await containersPromise}
		<Spinner />
	{:catch}
		<p>{$t("error.unexpected.title")}</p>
	{/await}
</Card>
