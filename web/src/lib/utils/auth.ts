import type { TokenPayload } from "../../domain/jwt";

/**
 * Cookie name for the JWT token.
 */
const TOKEN_COOKIE_NAME = "token";

/**
 * Retrieves subject token.
 * @returns Token or `null` when the token is not found.
 */
export function getToken(): string | null {
	// Get individual name=value cookies.
	const cookies = document.cookie.split(";");

	for (const cookie of cookies) {
		const [cookieName, cookieValue] = cookie.split("=");
		if (TOKEN_COOKIE_NAME == cookieName.trim()) {
			return decodeURIComponent(cookieValue);
		}
	}

	return null;
}

/**
 * Decodes payload from token.
 * @param token JWT token.
 * @returns Decoded payload or `null` when payload is not found within the token.
 */
export function decodeTokenPayload(token: string): TokenPayload | null {
	const [, payloadBase64] = token.split(".");

	if (!payloadBase64) {
		return null;
	}

	const decodedPayloadStr = atob(payloadBase64);
	const payload = JSON.parse(decodedPayloadStr);

	return payload;
}

/**
 * Stores token in cookies.
 * @param token JWT token.
 * @param expirationTime JWT expiration time.
 */
export function storeToken(token: string, expirationTime: number) {
	const expireTimeInMs = expirationTime * 1000;
	const expireDate = new Date(expireTimeInMs).toUTCString();

	document.cookie = `token=${token}; Path=/; Expires=${expireDate}; SameSite=Strict; Secure`;
}

/**
 * Clears token in cookies.
 */
export function clearToken() {
	document.cookie = `token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
}

/**
 * Indicates whether the employee is viewing himself/herself.
 * @param employeeId The employee ID.
 */
export function isViewingSelf(employeeId: string): boolean {
	const token = getToken();
	if (!token) {
		throw new Error("Failed to retrieve employee token");
	}

	const tokenPayload = decodeTokenPayload(token);
	if (!tokenPayload) {
		throw new Error("Failed to decode token payload");
	}

	return employeeId === tokenPayload.sub;
}
