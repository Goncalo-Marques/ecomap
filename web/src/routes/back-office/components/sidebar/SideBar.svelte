<script lang="ts">
	import { page } from "$app/stores";
	import Icon from "$lib/components/Icon.svelte";
	import { type LocaleTextID, t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../../constants/routes";
	import Support from "./Support.svelte";

	type MenuItem = {
		href: BackOfficeRoutes;
		icon: string;
		localeTextId: LocaleTextID;
	};

	const menuItems: MenuItem[] = [
		{
			href: BackOfficeRoutes.DASHBOARD,
			icon: "dashboard",
			localeTextId: "dashboard",
		},
		{
			href: BackOfficeRoutes.MAP,
			icon: "map",
			localeTextId: "map",
		},
		{
			href: BackOfficeRoutes.ROUTES,
			icon: "route",
			localeTextId: "routes",
		},
		{
			href: BackOfficeRoutes.CONTAINERS,
			icon: "delete",
			localeTextId: "containers",
		},
		{
			href: BackOfficeRoutes.TRUCKS,
			icon: "local_shipping",
			localeTextId: "trucks",
		},
		{
			href: BackOfficeRoutes.WAREHOUSES,
			icon: "warehouse",
			localeTextId: "warehouses",
		},
		{
			href: BackOfficeRoutes.LANDFILLS,
			icon: "landslide",
			localeTextId: "landfills",
		},
		{
			href: BackOfficeRoutes.EMPLOYEES,
			icon: "group",
			localeTextId: "employees",
		},
	];
</script>

<nav
	class="flex w-52 flex-col justify-between border-r border-gray-300 bg-white p-4"
>
	<ul class="flex flex-col justify-center gap-4">
		{#each menuItems as menuItem}
			<li class="block flex-1">
				<a
					class="flex w-full justify-items-center gap-2 rounded px-3 py-2 transition-colors hover:bg-green-50 hover:text-green-800 active:bg-green-100 active:text-green-800 data-[selected=true]:bg-green-100 data-[selected=true]:font-semibold data-[selected=true]:text-green-800"
					href={menuItem.href}
					data-selected={$page.url.pathname === menuItem.href ||
						$page.url.pathname.startsWith(`${menuItem.href}/`)}
				>
					<Icon name={menuItem.icon} />
					<span>{$t(menuItem.localeTextId)}</span>
				</a>
			</li>
		{/each}
	</ul>
	<Support />
</nav>
