<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../../../../lib/components/Icon.svelte";
	import url from "../../../../lib/stores/url";
	import { t, type LocaleTextID } from "../../../../lib/utils/i8n";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import Support from "./Support.svelte";

	type MenuItem = {
		route: BackOfficeRoutes;
		icon: string;
		localeTextId: LocaleTextID;
	};

	const menuItems: MenuItem[] = [
		{
			route: BackOfficeRoutes.DASHBOARD,
			icon: "dashboard",
			localeTextId: "dashboard",
		},
		{
			route: BackOfficeRoutes.MAP,
			icon: "map",
			localeTextId: "map",
		},
		{
			route: BackOfficeRoutes.ROUTES,
			icon: "route",
			localeTextId: "routes",
		},
		{
			route: BackOfficeRoutes.CONTAINERS,
			icon: "delete",
			localeTextId: "containers",
		},
		{
			route: BackOfficeRoutes.TRUCKS,
			icon: "local_shipping",
			localeTextId: "trucks",
		},
		{
			route: BackOfficeRoutes.WAREHOUSES,
			icon: "warehouse",
			localeTextId: "warehouses",
		},
		{
			route: BackOfficeRoutes.LANDFILLS,
			icon: "landslide",
			localeTextId: "landfills",
		},
		{
			route: BackOfficeRoutes.EMPLOYEES,
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
				<Link
					class="flex w-full justify-items-center gap-2 rounded px-3 py-2 transition-colors hover:bg-green-50 hover:text-green-800 active:bg-green-100 active:text-green-800 data-[selected=true]:bg-green-100 data-[selected=true]:font-semibold data-[selected=true]:text-green-800"
					to={menuItem.route}
					data-selected={$url.pathname === menuItem.route ||
						$url.pathname.startsWith(`${menuItem.route}/`)}
				>
					<Icon name={menuItem.icon} />
					<span>{$t(menuItem.localeTextId)}</span>
				</Link>
			</li>
		{/each}
	</ul>
	<Support />
</nav>
