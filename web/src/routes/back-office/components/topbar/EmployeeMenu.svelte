<script lang="ts">
	import { goto } from "$app/navigation";
	import type { Employee } from "$domain/employees";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Divider from "$lib/components/Divider.svelte";
	import Icon from "$lib/components/Icon.svelte";
	import Menu from "$lib/components/Menu.svelte";
	import MenuItem from "$lib/components/MenuItem.svelte";
	import Popover from "$lib/components/Popover.svelte";
	import { BackOfficeRoutes, CommonRoutes } from "$lib/constants/routes";
	import { clearToken, decodeTokenPayload, getToken } from "$lib/utils/auth";
	import { t } from "$lib/utils/i8n";

	/**
	 * The employee popover element.
	 */
	let popover: HTMLElement;

	/**
	 * Signs the employee out.
	 */
	function signOut() {
		clearToken();
		goto(CommonRoutes.SIGN_IN);
	}

	/**
	 * Retrieves the information about the employee who is currently signed in.
	 */
	async function getSignedInEmployee(): Promise<Employee> {
		const token = getToken();
		if (!token) {
			throw new Error("Failed to retrieve employee token");
		}

		const tokenPayload = decodeTokenPayload(token);
		if (!tokenPayload) {
			throw new Error("Failed to decode token payload");
		}

		const res = await ecomapHttpClient.GET("/employees/{employeeId}", {
			params: { path: { employeeId: tokenPayload.sub } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve employee details");
		}

		return res.data;
	}

	const employeePromise = getSignedInEmployee();
</script>

{#await employeePromise then employee}
	<Popover bind:popover id="employee-popover" align="right">
		<div slot="trigger" class="flex items-center justify-end gap-2">
			<div class="flex items-center gap-1">
				<Icon name="account_circle" />
				<span>{employee.firstName} {employee.lastName}</span>
			</div>
			<Icon name="expand_more" size="small" />
		</div>

		<Menu>
			<MenuItem>
				<a
					href={`${BackOfficeRoutes.EMPLOYEES}/${employee.id}`}
					class="contents"
					on:click={() => popover.hidePopover()}
				>
					<Icon name="person" size="small" />
					{$t("myInformation")}
				</a>
			</MenuItem>
			<Divider />
			<MenuItem>
				<button class="flex w-full items-center gap-2" on:click={signOut}>
					<Icon name="logout" size="small" />
					{$t("signOut.button")}
				</button>
			</MenuItem>
		</Menu>
	</Popover>
{:catch}
	{signOut()}
{/await}
