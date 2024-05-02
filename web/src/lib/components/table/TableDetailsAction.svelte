<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../Icon.svelte";

	/**
	 * Resource ID.
	 */
	export let id: string;
</script>

<Link style="display: contents;" to={id}>
	<Icon name="chevron_right" />
</Link>
