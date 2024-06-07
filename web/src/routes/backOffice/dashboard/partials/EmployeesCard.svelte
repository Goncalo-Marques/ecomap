<script lang="ts">
	import type { Employee } from "../../../../domain/employees";
	import Card from "../../components/Card.svelte";
</script>

<Card>
	<h2>Colaboradores ativos</h2>
	<p>
		<slot />
	</p>
</Card>
