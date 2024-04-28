/**
 * Debounces a given callback with a given time.
 * @template TCallback Type of the callback function.
 * @param callback Callback after the debounce.
 * @param waitFor Number of milliseconds to wait for the debounce.
 * @returns Debounced callback.
 */
function debounce<
	TCallback extends (...args: Parameters<TCallback>) => ReturnType<TCallback>,
>(callback: TCallback, waitFor: number) {
	let timeout: NodeJS.Timeout;

	return (...args: Parameters<TCallback>) => {
		clearTimeout(timeout);
		timeout = setTimeout(() => callback(...args), waitFor);
	};
}

export default debounce;
