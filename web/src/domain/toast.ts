/**
 * Represents the options available for the toast.
 */
export interface ToastOptions {
	/**
	 * The toast type.
	 */
	type: "success" | "error";

	/**
	 * The toast title.
	 */
	title: string;

	/**
	 * The toast description.
	 */
	description?: string;
}

/**
 * Represents the toast context used to interact with the toast.
 */
export interface ToastContext {
	/**
	 * Reveals a toast with a given set of options.
	 * @param options Toast options.
	 */
	show(options: ToastOptions): void;
}
