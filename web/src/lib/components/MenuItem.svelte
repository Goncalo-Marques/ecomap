<li>
	<slot />
</li>

<style>
	li {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 0.75rem;
		text-wrap: nowrap;
		transition-property: var(--transition-colors);
		transition-duration: var(--transition-duration);
		transition-timing-function: var(--transition-timing-function);

		& > * {
			flex: 1;
		}

		&:hover {
			color: var(--green-800);
			background-color: var(--green-50);
		}

		&:active {
			color: var(--green-800);
			background-color: var(--green-100);
		}
	}
</style>
