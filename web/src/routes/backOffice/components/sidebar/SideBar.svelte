<script lang="ts">
	import { Link } from "svelte-routing";
	import Icon from "../../../../lib/components/Icon.svelte";
	import url from "../../../../lib/stores/url";
	import { t } from "../../../../lib/utils/i8n";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import LocaleSelect from "./LocaleSelect.svelte";
</script>

<nav>
	<ul>
		<li>
			<Link
				to={BackOfficeRoutes.DASHBOARD}
				data-selected={$url.pathname === BackOfficeRoutes.DASHBOARD ||
					$url.pathname.startsWith(`${BackOfficeRoutes.DASHBOARD}/`)}
			>
				<Icon name="dashboard" />
				<span>{$t("dashboard")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.MAP}
				data-selected={$url.pathname === BackOfficeRoutes.MAP ||
					$url.pathname.startsWith(`${BackOfficeRoutes.MAP}/`)}
			>
				<Icon name="map" />
				<span>{$t("map")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.ROUTES}
				data-selected={$url.pathname === BackOfficeRoutes.ROUTES ||
					$url.pathname.startsWith(`${BackOfficeRoutes.ROUTES}/`)}
			>
				<Icon name="route" />
				<span>{$t("routes")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.CONTAINERS}
				data-selected={$url.pathname === BackOfficeRoutes.CONTAINERS ||
					$url.pathname.startsWith(`${BackOfficeRoutes.CONTAINERS}/`)}
			>
				<Icon name="delete" />
				<span>{$t("containers")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.TRUCKS}
				data-selected={$url.pathname === BackOfficeRoutes.TRUCKS ||
					$url.pathname.startsWith(`${BackOfficeRoutes.TRUCKS}/`)}
			>
				<Icon name="local_shipping" />
				<span>{$t("trucks")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.WAREHOUSES}
				data-selected={$url.pathname === BackOfficeRoutes.WAREHOUSES ||
					$url.pathname.startsWith(`${BackOfficeRoutes.WAREHOUSES}/`)}
			>
				<Icon name="warehouse" />
				<span>{$t("warehouses")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.REPORTS}
				data-selected={$url.pathname === BackOfficeRoutes.REPORTS ||
					$url.pathname.startsWith(`${BackOfficeRoutes.REPORTS}/`)}
			>
				<Icon name="description" />
				<span>{$t("reports")}</span>
			</Link>
		</li>
		<li>
			<Link
				to={BackOfficeRoutes.EMPLOYEES}
				data-selected={$url.pathname === BackOfficeRoutes.EMPLOYEES ||
					$url.pathname.startsWith(`${BackOfficeRoutes.EMPLOYEES}/`)}
			>
				<Icon name="group" />
				<span>{$t("employees")}</span>
			</Link>
		</li>
	</ul>
	<LocaleSelect />
</nav>

<style>
	nav {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		width: 13rem;
		padding: 1rem;
		background-color: var(--white);
		border-right: 1px solid var(--gray-300);
		flex-shrink: 0;
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
