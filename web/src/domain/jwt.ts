import type { SubjectRole } from "./role";

/**
 * JWT token payload.
 */
export interface TokenPayload {
	/**
	 * Issuer claim.
	 * Identifies the principal that issued the JWT.
	 */
	iss: string;

	/**
	 * Subject claim.
	 * Identifies the principal that is the subject of the JWT.
	 */
	sub: string;

	/**
	 * Expiration time claim.
	 * Identifies the expiration time on or after which the JWT must not be accepted for processing.
	 */
	exp: number;

	/**
	 * Issued at claim.
	 * Identifies the time at which the JWT was issued.
	 */
	iat: number;

	/**
	 * Roles of the subject.
	 */
	roles: SubjectRole[];
}
