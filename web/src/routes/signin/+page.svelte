<script lang="ts">
	import { goto } from "$app/navigation";
	import type { TokenPayload } from "$domain/jwt";
	import { SubjectRole } from "$domain/role";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import FormControl from "$lib/components/FormControl.svelte";
	import Input from "$lib/components/Input.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { decodeTokenPayload, getToken, storeToken } from "$lib/utils/auth";
	import { t } from "$lib/utils/i8n";

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
			goto(BackOfficeRoutes.DASHBOARD, { replaceState: true });
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

	const roles = getUserRoles();
	if (roles.includes(SubjectRole.MANAGER)) {
		// Redirect to back office dashboard page if user is a manager.
		goto(BackOfficeRoutes.DASHBOARD, { replaceState: true });
	}
</script>

<main class="grid min-h-screen place-items-center">
	<article class="flex w-96 flex-col gap-8 rounded-3xl bg-white p-12 shadow-md">
		<h1 class="text-4xl font-semibold">{$t("signin.title")}</h1>

		<form
			novalidate
			method="post"
			class="flex flex-col gap-8"
			on:submit|preventDefault={handleSubmit}
		>
			<div class="flex flex-col gap-4">
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
					<p class="text-xs text-red-500">{responseErrorMessage}</p>
				{/if}
			</div>

			<Button size="large" type="submit">
				{$t("signin.button")}
			</Button>
		</form>
	</article>
</main>
