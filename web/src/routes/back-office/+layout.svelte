<script lang="ts">
	import Toast from "$lib/components/Toast.svelte";

	import SideBar from "./components/sidebar/SideBar.svelte";
	import TopBar from "./components/topbar/TopBar.svelte";
</script>

<Toast>
	<TopBar />
	<div class="flex h-[calc(100vh-var(--top-bar-height))]">
		<SideBar />
		<div class="grid flex-1 overflow-auto">
			<slot />
		</div>
	</div>
</Toast>
