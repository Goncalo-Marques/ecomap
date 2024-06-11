import "./app.css";
import App from "./App.svelte";
import { setupChart } from "./lib/utils/chart";

setupChart();

const app = new App({
	target: document.getElementById("app")!,
});

export default app;
