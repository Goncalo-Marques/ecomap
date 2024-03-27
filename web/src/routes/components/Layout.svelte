<script>
	import SideBar from "./sidebar/SideBar.svelte";
	import TopBar from "./topbar/TopBar.svelte";
</script>

<div class="main">
	<TopBar />
	<div class="container">
		<SideBar />
		<slot />
	</div>
</div>

<style>
	.main {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.container {
		flex: 1;
		display: flex;
	}
</style>
