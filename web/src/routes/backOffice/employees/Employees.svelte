<!-- TODO: Replace the current page with the actual employee data -->
<script lang="ts">
	import type { ComponentProps } from "svelte";
	import Button from "../../../lib/components/Button.svelte";
	import Table from "../../../lib/components/table/Table.svelte";
	import type { Columns, Sorting } from "../../../lib/components/table/types";
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import type { components } from "../../../../api/ecomap/http";

	type User = components["schemas"]["User"];

	let usersPromise: Promise<User[]>;
	let usersAmount = 0;

	const pageSize = 1;

	let sorting: Sorting<User> = {
		field: "firstName",
		direction: "asc",
	};

	const columns: Columns<User> = [
		{
			type: "accessor",
			field: "firstName",
			header: "Name",
			enableSorting: true,
			cell(firstName, row) {
				return `${firstName} ${row.lastName}`;
			},
		},
		{
			type: "accessor",
			field: "username",
			header: "Username",
			enableSorting: true,
			cell(username) {
				return username;
			},
		},
		{
			type: "accessor",
			field: "createdTime",
			header: "Created at",
			enableSorting: false,
			cell(createdTime) {
				return new Date(createdTime).toDateString();
			},
		},
		{
			type: "accessor",
			field: "modifiedTime",
			header: "Modified at",
			enableSorting: true,
			cell(modifiedTime) {
				return new Date(modifiedTime).toDateString();
			},
		},
		{
			type: "display",
			header: "",
			align: "center",
			size: 120,
			cell() {
				return {
					component: Button,
					props: {
						variant: "tertiary",
						size: "medium",
					} as ComponentProps<Button>,
					slotContent: "View details",
				};
			},
		},
	];

	let pageIndex = 0;

	async function fetchUsers(
		pageIndex: number,
		pageSize: number,
		sorting: Sorting<User>,
	): Promise<User[]> {
		const res = await ecomapHttpClient.GET("/users", {
			params: {
				query: {
					offset: pageIndex * pageSize,
					limit: pageSize,
					sort: sorting.field,
					order: sorting.direction,
				},
			},
		});

		if (res.error) {
			usersAmount = 0;
			return [];
		}

		usersAmount = res.data.total;
		return res.data.users;
	}

	function handleSortingChange(newSorting: Sorting<User>) {
		usersPromise = fetchUsers(pageIndex, pageSize, newSorting);
		sorting = newSorting;
	}

	async function handlePageChange(newPageIndex: number) {
		usersPromise = fetchUsers(newPageIndex, pageSize, sorting);
		pageIndex = newPageIndex;
	}

	usersPromise = fetchUsers(pageIndex, pageSize, sorting);
</script>

<main>
	<Table
		{columns}
		{sorting}
		rows={usersPromise}
		pagination={{
			name: "users",
			pageIndex,
			pageSize,
			total: usersAmount,
			onPageChange: handlePageChange,
		}}
		onSortingChange={handleSortingChange}
	/>
</main>

<style>
	main {
		flex: 1;
		padding: 2.5rem;
	}
</style>
