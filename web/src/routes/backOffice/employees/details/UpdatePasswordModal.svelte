<script lang="ts">
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import Input from "../../../../lib/components/Input.svelte";
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { Employee } from "../../../../domain/employees";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import FormModal from "../../../../lib/components/FormModal.svelte";
	import { decodeTokenPayload, storeToken } from "../../../../lib/utils/auth";

	/**
	 * The employee whose password needs to be changed.
	 */
	export let employee: Employee;

	/**
	 * Indicates if the modal is open.
	 */
	export let open: boolean;

	/**
	 * Callback fired when open state of the modal changes.
	 * @param open New open state modal.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * The form element.
	 */
	let form: HTMLFormElement | null = null;

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		currentPassword: "",
		newPassword: "",
		confirmPassword: "",
	};

	/**
	 * Validates the form and sets error messages on the form fields
	 * if they contain any errors.
	 * @param currentPasswordInput Current password input field.
	 * @param newPasswordInput New password input field.
	 * @param confirmPasswordInput Confirm password input field.
	 * @returns `true` if form is valid, `false` otherwise.
	 */
	function validateForm(
		currentPasswordInput: HTMLInputElement,
		newPasswordInput: HTMLInputElement,
		confirmPasswordInput: HTMLInputElement,
	) {
		if (currentPasswordInput.validity.valueMissing) {
			formErrorMessages.currentPassword = $t("error.valueMissing");
		} else {
			formErrorMessages.currentPassword = "";
		}

		if (newPasswordInput.validity.valueMissing) {
			formErrorMessages.newPassword = $t("error.valueMissing");
		} else {
			formErrorMessages.newPassword = "";
		}

		if (confirmPasswordInput.validity.valueMissing) {
			formErrorMessages.confirmPassword = $t("error.valueMissing");
		} else if (newPasswordInput.value !== confirmPasswordInput.value) {
			formErrorMessages.confirmPassword = $t(
				"employees.error.passwordMismatch",
			);
		} else {
			formErrorMessages.confirmPassword = "";
		}

		return (
			!formErrorMessages.currentPassword &&
			!formErrorMessages.newPassword &&
			!formErrorMessages.confirmPassword &&
			newPasswordInput.value === confirmPasswordInput.value
		);
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

		const currentPassword = formData.get("currentPassword") ?? "";
		const newPassword = formData.get("newPassword") ?? "";
		const confirmPassword = formData.get("confirmPassword") ?? "";

		// Check if fields are strings.
		if (
			typeof currentPassword !== "string" ||
			typeof newPassword !== "string" ||
			typeof confirmPassword !== "string"
		) {
			return;
		}

		const currentPasswordInput = form.elements.namedItem(
			"currentPassword",
		) as HTMLInputElement;
		const newPasswordInput = form.elements.namedItem(
			"newPassword",
		) as HTMLInputElement;
		const confirmPasswordInput = form.elements.namedItem(
			"confirmPassword",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				currentPasswordInput,
				newPasswordInput,
				confirmPasswordInput,
			)
		) {
			return;
		}

		const updatePasswordRes = await ecomapHttpClient.PUT(
			"/employees/password",
			{
				body: {
					username: employee.username,
					oldPassword: currentPassword,
					newPassword: newPassword,
				},
			},
		);

		if (updatePasswordRes.error) {
			switch (updatePasswordRes.error.code) {
				case "bad_request":
					toast.show({
						type: "error",
						title: $t(
							"employees.updatePassword.error.passwordConstraints.title",
						),
						description: $t(
							"employees.updatePassword.error.passwordConstraints.description",
						),
					});
					break;

				case "unauthorized":
					toast.show({
						type: "error",
						title: $t("employees.updatePassword.error.incorrectPassword.title"),
						description: $t(
							"employees.updatePassword.error.incorrectPassword.description",
						),
					});
					break;

				default:
					toast.show({
						type: "error",
						title: $t("error.unexpected.title"),
						description: $t("error.unexpected.description"),
					});
					break;
			}
			return;
		}

		// Sign in employee with the new password to retrieve the new JWT token.
		const signinRes = await ecomapHttpClient.POST("/employees/signin", {
			body: {
				username: employee.username,
				password: newPassword,
			},
		});

		if (signinRes.error) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});
			return;
		}

		// Retrieve the JWT token.
		const token = signinRes.data.token;

		const tokenPayload = decodeTokenPayload(token);
		if (!tokenPayload) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});
			return;
		}

		// Store token in cookies.
		storeToken(token, tokenPayload.exp);

		toast.show({
			type: "success",
			title: $t("employees.updatePassword.success"),
			description: undefined,
		});
		onOpenChange(false);
	}

	/**
	 * Handles cancel action.
	 */
	function handleCancel() {
		onOpenChange(false);
	}

	// Clear form when modal is closed.
	$: if (!open && form) {
		// Reset form to its initial state.
		form.reset();

		// Reset error messages on the inputs.
		formErrorMessages = {
			currentPassword: "",
			newPassword: "",
			confirmPassword: "",
		};
	}
</script>

<FormModal
	{open}
	{onOpenChange}
	bind:form
	gutters
	title={$t("employees.updatePassword.title")}
	onSubmit={handleSubmit}
>
	<div class="content">
		<FormControl
			label={$t("employees.updatePassword.currentPassword.label")}
			error={!!formErrorMessages.currentPassword}
			helperText={formErrorMessages.currentPassword}
		>
			<Input
				required
				type="password"
				name="currentPassword"
				placeholder={$t("employees.updatePassword.currentPassword.placeholder")}
				error={!!formErrorMessages.currentPassword}
			/>
		</FormControl>
		<FormControl
			label={$t("employees.updatePassword.newPassword.label")}
			error={!!formErrorMessages.newPassword}
			helperText={formErrorMessages.newPassword}
			title={$t("employees.passwordConstraints")}
		>
			<Input
				required
				type="password"
				name="newPassword"
				placeholder={$t("employees.updatePassword.newPassword.placeholder")}
				error={!!formErrorMessages.newPassword}
			/>
		</FormControl>
		<FormControl
			label={$t("employees.updatePassword.confirmPassword.label")}
			error={!!formErrorMessages.confirmPassword}
			helperText={formErrorMessages.confirmPassword}
		>
			<Input
				required
				type="password"
				name="confirmPassword"
				placeholder={$t("employees.updatePassword.confirmPassword.placeholder")}
				error={!!formErrorMessages.confirmPassword}
			/>
		</FormControl>
	</div>
	<svelte:fragment slot="actions">
		<Button type="button" variant="tertiary" onClick={handleCancel}>
			{$t("cancel")}
		</Button>
		<Button type="submit">{$t("confirm")}</Button>
	</svelte:fragment>
</FormModal>

<style>
	.content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		width: 40rem;
	}
</style>
