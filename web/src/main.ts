import {
	ArcElement,
	BarController,
	BarElement,
	CategoryScale,
	Chart,
	DoughnutController,
	Legend,
	LineController,
	LineElement,
	LinearScale,
	PointElement,
	Tooltip,
} from "chart.js";
import "./app.css";
import App from "./App.svelte";
import { getCssVariable } from "./lib/utils/cssVars";

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
);
Chart.defaults.font.family = "Inter, sans-serif";
Chart.defaults.font.lineHeight = 1.5;
Chart.defaults.color = getCssVariable("--gray-900");

const app = new App({
	target: document.getElementById("app")!,
});

export default app;
