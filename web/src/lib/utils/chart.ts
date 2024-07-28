import {
	ArcElement,
	BarController,
	BarElement,
	CategoryScale,
	Chart,
	DoughnutController,
	Filler,
	Legend,
	LinearScale,
	LineController,
	LineElement,
	PointElement,
	Tooltip,
} from "chart.js";

import { getCssVariable } from "./cssVars";

/**
 * Set up Chart.js with the required components for the application and a
 * custom theme.
 */
export function setupChart() {
	Chart.register(
		LineElement,
		LineController,
		LinearScale,
		CategoryScale,
		PointElement,
		Tooltip,
		Legend,
		DoughnutController,
		ArcElement,
		BarElement,
		BarController,
		Filler,
	);

	Chart.defaults.font.family = "Inter, sans-serif";
	Chart.defaults.font.lineHeight = 1.5;
	Chart.defaults.color = getCssVariable("--gray-900");
}
