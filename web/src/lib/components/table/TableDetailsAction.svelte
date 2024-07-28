<script lang="ts">
	import Icon from "../Icon.svelte";

	/**
	 * Resource href.
	 */
	export let href: string;
</script>

<a class="contents" {href}>
	<Icon name="chevron_right" />
</a>
