<script lang="ts">
	let inputValue: string;
	let promise: Promise<{ id: string; name: string; dateOfBirth: string }>;

	async function getEmployee(employeeId: string) {
		const res = await fetch(`/api/employees/${employeeId}`);
		const body = await res.json();

		if (res.ok) {
			return body;
		} else {
			throw new Error(body.message);
		}
	}

	function handleClick() {
		promise = getEmployee(inputValue);
	}
</script>

<p>
	<label for="employeeId">Employee ID</label>
	<input id="employeeId" type="text" bind:value={inputValue} />
</p>

<button on:click={handleClick}>Search</button>

{#await promise}
	<p>Loading...</p>
{:then employee}
	{#if employee}
		<p>{JSON.stringify(employee, null, 2)}</p>
	{/if}
{:catch error}
	<p style="color: red">{error.message}</p>
{/await}

<style>
	button {
		width: fit-content;
	}
</style>
