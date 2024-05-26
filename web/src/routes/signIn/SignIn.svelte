<script lang="ts">
	import { onMount } from "svelte";
	import { navigate } from "svelte-routing";
	import { BackOfficeRoutes } from "../constants/routes";
	import Button from "../../lib/components/Button.svelte";
	import Input from "../../lib/components/Input.svelte";
	import { t } from "../../lib/utils/i8n";
	import {
		decodeTokenPayload,
		getToken,
		storeToken,
	} from "../../lib/utils/auth";
	import ecomapHttpClient from "../../lib/clients/ecomap/http";
	import { SubjectRole } from "../../domain/role";
	import type { TokenPayload } from "../../domain/jwt";
	import FormControl from "../../lib/components/FormControl.svelte";

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
			formErrorMessages.username = $t("error.valueMissing");
		} else {
			formErrorMessages.username = "";
		}

		if (!password) {
			formErrorMessages.password = $t("error.valueMissing");
		} else {
			formErrorMessages.password = "";
		}
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	async function handleSubmit(
		e: Event & { currentTarget: EventTarget & HTMLFormElement },
	) {
		const form = e.currentTarget;
		const formData = new FormData(form);
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

		const res = await ecomapHttpClient.POST("/employees/signin", {
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

		const { token } = res.data;

		const payload = decodeTokenPayload(token);
		if (!payload) {
			responseErrorMessage = $t("error.unexpected");
			return;
		}

		if (payload.roles.includes(SubjectRole.MANAGER)) {
			storeToken(token, payload.exp);

			// Redirect to back office dashboard page if user is a manager.
			navigate(BackOfficeRoutes.DASHBOARD, { replace: true });
		} else {
			responseErrorMessage = $t("error.unexpected");
		}
	}

	/**
	 * Retrieves user roles.
	 * @returns User roles.
	 */
	function getUserRoles(): TokenPayload["roles"] {
		const token = getToken();
		if (!token) {
			return [];
		}

		const payload = decodeTokenPayload(token);
		if (!payload) {
			return [];
		}

		return payload.roles;
	}

	onMount(() => {
		const roles = getUserRoles();
		if (roles.includes(SubjectRole.MANAGER)) {
			// Redirect to back office dashboard page if user is a manager.
			navigate(BackOfficeRoutes.DASHBOARD, { replace: true });
		}
	});
</script>

<main>
	<article>
		<h1>{$t("signin.title")}</h1>

		<form novalidate method="post" on:submit|preventDefault={handleSubmit}>
			<div class="container">
				<FormControl
					error={!!formErrorMessages.username}
					label={$t("signin.username.label")}
					helperText={formErrorMessages.username}
				>
					<Input
						required
						type="text"
						name="username"
						error={!!formErrorMessages.username}
						autocomplete="off"
						placeholder={$t("signin.username.placeholder")}
					/>
				</FormControl>

				<FormControl
					error={!!formErrorMessages.password}
					label={$t("signin.password.label")}
					helperText={formErrorMessages.password}
				>
					<Input
						required
						type="password"
						name="password"
						error={!!formErrorMessages.password}
						placeholder={$t("signin.password.placeholder")}
					/>
				</FormControl>

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
