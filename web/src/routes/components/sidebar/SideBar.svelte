<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../../../lib/components/Icon.svelte";
	import url from "../../../lib/utils/url";
	import { t } from "../../../lib/utils/i8n";
	import { AppRoutes } from "../../constants/appRoutes";
</script>

<nav>
	<ul>
		<li>
			<Link
				to={AppRoutes.DASHBOARD}
				data-selected={$url.pathname === AppRoutes.DASHBOARD ||
					$url.pathname.startsWith(`${AppRoutes.DASHBOARD}/`)}
			>
				<Icon name="dashboard" />
				<span>{$t("sidebar.dashboard")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.MAP}
				data-selected={$url.pathname === AppRoutes.MAP ||
					$url.pathname.startsWith(`${AppRoutes.MAP}/`)}
			>
				<Icon name="map" />
				<span>{$t("sidebar.map")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.ROUTES}
				data-selected={$url.pathname === AppRoutes.ROUTES ||
					$url.pathname.startsWith(`${AppRoutes.ROUTES}/`)}
			>
				<Icon name="route" />
				<span>{$t("sidebar.routes")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.CONTAINERS}
				data-selected={$url.pathname === AppRoutes.CONTAINERS ||
					$url.pathname.startsWith(`${AppRoutes.CONTAINERS}/`)}
			>
				<Icon name="delete" />
				<span>{$t("sidebar.containers")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.WAREHOUSES}
				data-selected={$url.pathname === AppRoutes.WAREHOUSES ||
					$url.pathname.startsWith(`${AppRoutes.WAREHOUSES}/`)}
			>
				<Icon name="warehouse" />
				<span>{$t("sidebar.warehouses")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.TRUCKS}
				data-selected={$url.pathname === AppRoutes.TRUCKS ||
					$url.pathname.startsWith(`${AppRoutes.TRUCKS}/`)}
			>
				<Icon name="local_shipping" />
				<span>{$t("sidebar.trucks")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.REPORTS}
				data-selected={$url.pathname === AppRoutes.REPORTS ||
					$url.pathname.startsWith(`${AppRoutes.REPORTS}/`)}
			>
				<Icon name="description" />
				<span>{$t("sidebar.reports")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={AppRoutes.EMPLOYEES}
				data-selected={$url.pathname === AppRoutes.EMPLOYEES ||
					$url.pathname.startsWith(`${AppRoutes.EMPLOYEES}/`)}
			>
				<Icon name="group" />
				<span>{$t("sidebar.employees")}</span>
			</Link>
		</li>
	</ul>
</nav>

<style>
	nav {
		width: 13rem;
		padding: 1rem;
		background-color: var(--white);
		border-right: 1px solid var(--gray-300);
	}

	ul {
		display: flex;
		flex-direction: column;
		justify-items: center;
		gap: 1rem;
	}

	li {
		flex: 1;
		display: block;
	}

	nav :global(a) {
		display: flex;
		justify-items: center;
		gap: 0.5rem;
		border-radius: 0.25rem;
		width: 100%;
		padding: 0.5rem 0.75rem;
		transition-property: var(--transition-colors);
		transition-duration: var(--transition-duration);
		transition-timing-function: var(--transition-timing-function);

		&:hover {
			background-color: var(--green-50);
			color: var(--green-800);
		}

		&:active {
			background-color: var(--green-100);
			color: var(--green-800);
		}

		&[data-selected="true"] {
			background-color: var(--green-100);
			color: var(--green-800);
			font-weight: 600;
		}
	}
</style>
