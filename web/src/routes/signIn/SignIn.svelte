<script lang="ts">
	import { navigate } from "svelte-routing";
	import { BackOfficeRoutes } from "../constants/routes";
	import Button from "../../lib/components/Button.svelte";
	import Input from "../../lib/components/Input.svelte";
	import httpClient from "../../lib/utils/httpClient";
	import { t } from "../../lib/utils/i8n";
	import { storeToken } from "../../lib/utils/auth";

	/**
	 * Error message displayed after an error occurs with the server.
	 */
	let responseErrorMessage: string;

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		username: "",
		password: "",
	};

	/**
	 * Validates the form and sets error messages on the form fields
	 * if they contain any errors.
	 * @param username Username field value.
	 * @param password Password field value.
	 */
	function validateForm(username: string, password: string) {
		if (!username) {
			formErrorMessages.username = $t("error.requiredField");
		} else {
			formErrorMessages.username = "";
		}

		if (!password) {
			formErrorMessages.password = $t("error.requiredField");
		} else {
			formErrorMessages.password = "";
		}
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	async function handleSubmit(e: SubmitEvent) {
		const formData = new FormData(e.currentTarget as HTMLFormElement);
		const username = formData.get("username");
		const password = formData.get("password");

		// Check if username and password are both strings.
		if (typeof username !== "string" || typeof password !== "string") {
			return;
		}

		validateForm(username, password);

		// Check if either fields are not filled to prevent making a server request.
		if (!username || !password) {
			return;
		}

		const res = await httpClient.POST("/employees/signin", {
			body: {
				username,
				password,
			},
		});

		if (res.error) {
			const { code } = res.error;

			if (code === "unauthorized") {
				responseErrorMessage = $t("signin.error.incorrectCredentials");
			} else {
				responseErrorMessage = $t("error.unexpected");
			}

			return;
		}

		try {
			storeToken(res.data.token);
		} catch {
			responseErrorMessage = $t("error.unexpected");
			return;
		}

		navigate(BackOfficeRoutes.DASHBOARD);
	}
</script>

<main>
	<article>
		<h1>{$t("signin.title")}</h1>

		<form method="post" on:submit|preventDefault={handleSubmit}>
			<div class="container">
				<Input
					type="text"
					name="username"
					error={!!formErrorMessages.username}
					helperText={formErrorMessages.username}
					autocomplete="off"
					label={$t("signin.username.label")}
					placeholder={$t("signin.username.placeholder")}
				/>

				<Input
					type="password"
					name="password"
					error={!!formErrorMessages.password}
					helperText={formErrorMessages.password}
					label={$t("signin.password.label")}
					placeholder={$t("signin.password.placeholder")}
				/>

				{#if responseErrorMessage}
					<p class="error">{responseErrorMessage}</p>
				{/if}
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

	.error {
		font: var(--text-xs-regular);
		color: var(--red-500);
	}
</style>
