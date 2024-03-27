<script lang="ts">
	import Button from "../../lib/components/Button.svelte";
	import Input from "../../lib/components/Input.svelte";
	import { t } from "../../lib/utils/i8n";
</script>

<main>
	<article>
		<h1>{$t("signin.title")}</h1>

		<form method="post" on:submit|preventDefault>
			<div class="container">
				<Input
					type="text"
					name="username"
					autocomplete="off"
					label={$t("signin.username.label")}
					placeholder={$t("signin.username.placeholder")}
				/>

				<Input
					type="password"
					name="password"
					label={$t("signin.password.label")}
					placeholder={$t("signin.password.placeholder")}
				/>
			</div>

			<Button size="large" type="submit">
				{$t("signin.button")}
			</Button>
		</form>
	</article>
</main>

<style>
	main {
		min-height: 100vh;
		display: grid;
		place-items: center;
	}

	article {
		background-color: var(--white);
		border-radius: 1.5rem;
		padding: 3rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
		width: 24rem;
		box-shadow: var(--shadow-md);
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.container {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
</style>
