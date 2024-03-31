import type { TokenPayload } from "../../domain/jwt";

/**
 * Cookie name for the JWT token.
 */
const TOKEN_COOKIE_NAME = "token" as const;

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
 * Indicates if subject is authenticated.
 * @returns `true` when subject is authenticated. Otherwise, returns `false`.
 */
export function isAuthenticated(): boolean {
	const token = getToken();

	return !!token;
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
 * @throws {Error} Invalid token payload.
 */
export function storeToken(token: string) {
	const payload = decodeTokenPayload(token);
	if (!payload) {
		throw new Error("Couldn't decode token payload");
	}

	const expireTimeInMs = payload.exp * 1000;
	const expireDate = new Date(expireTimeInMs);

	document.cookie = `token=${token}; Path=/; Expires=${expireDate}; SameSite=Strict; Secure`;
}
