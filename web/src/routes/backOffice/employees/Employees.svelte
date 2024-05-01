<!-- TODO: Replace the current page with the actual employee data -->
<script lang="ts">
	import { onMount, type ComponentProps } from "svelte";
	import Button from "../../../lib/components/Button.svelte";
	import Table from "../../../lib/components/table/Table.svelte";
	import type {
		Columns,
		SortingDirection,
	} from "../../../lib/components/table/types";
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import type { components } from "../../../../api/ecomap/http";

	type User = components["schemas"]["User"];
	type UserSortableFields = NonNullable<
		components["parameters"]["UserSortQueryParam"]
	>;

	let users: User[];
	let usersAmount = 0;

	let loading = false;

	const pageSize = 1;

	let sortingField: UserSortableFields = "firstName";
	let sortingOrder: SortingDirection = "asc";

	const columns: Columns<User> = [
		{
			type: "accessor",
			field: "firstName",
			header: "Name",
			enableSorting: true,
			enableFiltering: false,
			cell(firstName, row) {
				return `${firstName} ${row.lastName}`;
			},
		},
		{
			type: "accessor",
			field: "username",
			header: "Username",
			enableSorting: true,
			enableFiltering: false,
			cell(username) {
				return username;
			},
		},
		{
			type: "accessor",
			field: "createdAt",
			header: "Created at",
			enableSorting: false,
			enableFiltering: false,
			cell(createdAt) {
				return new Date(createdAt).toDateString();
			},
		},
		{
			type: "accessor",
			field: "modifiedAt",
			header: "Modified at",
			enableSorting: true,
			enableFiltering: false,
			cell(modifiedAt) {
				return new Date(modifiedAt).toDateString();
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
		sortingField: UserSortableFields,
		sortingOrder: SortingDirection,
	) {
		loading = true;

		const res = await ecomapHttpClient.GET("/users", {
			params: {
				query: {
					offset: pageIndex * pageSize,
					limit: pageSize,
					sort: sortingField,
					order: sortingOrder,
				},
			},
		});

		loading = false;

		if (res.error) {
			usersAmount = 0;
			users = [];
			return;
		}

		usersAmount = res.data.total;
		users = res.data.users;
	}

	function handleSortingChange(
		field: UserSortableFields,
		order: SortingDirection,
	) {
		fetchUsers(pageIndex, pageSize, field, order);
		sortingField = field;
		sortingOrder = order;
	}

	async function handlePageChange(newPageIndex: number) {
		fetchUsers(newPageIndex, pageSize, sortingField, sortingOrder);
		pageIndex = newPageIndex;
	}

	onMount(() => {
		fetchUsers(pageIndex, pageSize, sortingField, sortingOrder);
	});
</script>

<main>
	<Table
		{columns}
		{sortingField}
		{sortingOrder}
		{loading}
		rows={users}
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
