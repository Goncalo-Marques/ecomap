<dl>
	<slot />
</dl>

<style>
	dl {
		--column-count: 2;
		--column-gap: 0.5rem;
		--column-min-width: 13rem;

		display: grid;
		grid-template-columns: repeat(
			auto-fill,
			minmax(
				max(
					var(--column-min-width),
					(100% - var(--column-gap)) / var(--column-count)
				),
				1fr
			)
		);
		gap: var(--column-gap);
	}
</style>
