<script lang="ts">
	import {
		addVectorLayer,
		setStandard,
		setHumanitarian,
	} from "../utils/mapStore";
</script>

<main>
	<div class="side">
		<button
			on:click={() => {
				addVectorLayer("http://localhost:8000/concelhos.geojson");
			}}
		>
			Poly
		</button>
		<button
			on:click={() => {
				addVectorLayer("http://localhost:8000/empresas_baixo_vouga.geojson");
			}}
		>
			Point
		</button>

		<h2>Layers</h2>

		<div>
			<input
				type="radio"
				name="baseLayer"
				id="osmStandard"
				value="osmStandard"
				checked
				on:click={setStandard}
			/><label for="osmStandard">osmStandard</label>
		</div>
		<div>
			<input
				type="radio"
				name="baseLayer"
				id="osmHumanitarian"
				value="osmHumanitarian"
				on:click={setHumanitarian}
			/><label for="osmHumanitarian">osmHumanitarian</label>
		</div>
	</div>
	<div class="view">
		<slot />
	</div>
</main>

<style>
	* {
		box-sizing: border-box;
	}

	main {
		display: flex;
	}

	.side {
		display: flex;
		flex-direction: column;
		gap: 10px;
		float: left;

		width: 200px;
		height: 100vh;
		padding: 10px;
	}

	.view {
		float: none;
		width: 100%;
	}
</style>
