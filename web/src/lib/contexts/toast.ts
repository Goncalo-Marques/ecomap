import { getContext, setContext } from "svelte";
import type { ToastContext } from "../../domain/toast";

const TOAST_CONTEXT_KEY = "toast";

/**
 * Retrieves the toast context.
 * @returns Toast context.
 */
export function getToastContext(): ToastContext {
	return getContext<ToastContext>(TOAST_CONTEXT_KEY);
}

/**
 * Sets the toast context.
 * @returns Toast context.
 */
export function setToastContext(context: ToastContext): ToastContext {
	return setContext<ToastContext>(TOAST_CONTEXT_KEY, context);
}
