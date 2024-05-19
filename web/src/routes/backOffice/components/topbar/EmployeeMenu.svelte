<script lang="ts">
	import { Link, navigate } from "svelte-routing";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import Icon from "../../../../lib/components/Icon.svelte";
	import Popover from "../../../../lib/components/Popover.svelte";
	import {
		clearToken,
		decodeTokenPayload,
		getToken,
	} from "../../../../lib/utils/auth";
	import { BackOfficeRoutes, CommonRoutes } from "../../../constants/routes";
	import type { Employee } from "../../../../domain/employees";
	import MenuItem from "../../../../lib/components/MenuItem.svelte";
	import Menu from "../../../../lib/components/Menu.svelte";
	import Divider from "../../../../lib/components/Divider.svelte";
	import { t } from "../../../../lib/utils/i8n";

	/**
	 * The employee popover element.
	 */
	let popover: HTMLElement;

	/**
	 * Signs the employee out.
	 */
	function signOut() {
		clearToken();
		navigate(CommonRoutes.SIGN_IN);
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
		<div slot="trigger" class="employee-menu-trigger">
			<div class="avatar">
				<Icon name="account_circle" />
				<span>{employee.firstName} {employee.lastName}</span>
			</div>
			<Icon name="expand_more" size="small" />
		</div>

		<Menu>
			<MenuItem>
				<Link
					to={`${BackOfficeRoutes.EMPLOYEES}/${employee.id}`}
					style="display:contents"
					on:click={() => popover.hidePopover()}
				>
					<Icon name="person" size="small" />
					{$t("myInformation")}
				</Link>
			</MenuItem>
			<Divider />
			<MenuItem>
				<button class="signOut" on:click={signOut}>
					<Icon name="logout" size="small" />
					{$t("signOut.button")}
				</button>
			</MenuItem>
		</Menu>
	</Popover>
{:catch}
	{signOut()}
{/await}

<style>
	.signOut {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
	}

	.avatar {
		display: flex;
		align-items: center;
		gap: 0.25rem;
	}

	.employee-menu-trigger {
		display: flex;
		gap: 0.5rem;
		justify-content: flex-end;
		align-items: center;
	}
</style>
