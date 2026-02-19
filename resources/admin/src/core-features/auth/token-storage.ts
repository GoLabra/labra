'use client'

export const ACCESS_TOKEN_STORAGE_KEY = 'accessToken';
export const REFRESH_TOKEN_STORAGE_KEY = 'refreshToken';

type TokenPairPayload = {
    token?: string;
    refresh_token?: string;
};

let refreshRequestInFlight: Promise<string | null> | null = null;

const hasWindow = (): boolean => typeof window !== 'undefined';

const readStorage = (key: string): string | null => {
    if (!hasWindow()) {
        return null;
    }

    return window.localStorage.getItem(key);
};

const writeStorage = (key: string, value: string): void => {
    if (!hasWindow()) {
        return;
    }

    window.localStorage.setItem(key, value);
};

const removeStorage = (key: string): void => {
    if (!hasWindow()) {
        return;
    }

    window.localStorage.removeItem(key);
};

export const getAccessToken = (): string | null => readStorage(ACCESS_TOKEN_STORAGE_KEY);
export const getRefreshToken = (): string | null => readStorage(REFRESH_TOKEN_STORAGE_KEY);

export const setAccessToken = (token: string): void => writeStorage(ACCESS_TOKEN_STORAGE_KEY, token);

export const setAuthTokens = (token: string, refreshToken: string): void => {
    writeStorage(ACCESS_TOKEN_STORAGE_KEY, token);
    writeStorage(REFRESH_TOKEN_STORAGE_KEY, refreshToken);
};

export const clearAuthTokens = (): void => {
    removeStorage(ACCESS_TOKEN_STORAGE_KEY);
    removeStorage(REFRESH_TOKEN_STORAGE_KEY);
};

const extractTokenPair = (payload: TokenPairPayload): { token: string; refreshToken: string } | null => {
    const token = payload?.token?.trim();
    const refreshToken = payload?.refresh_token?.trim();

    if (!token || !refreshToken) {
        return null;
    }

    return { token, refreshToken };
};

const fetchRefreshedAccessToken = async (apiBaseUrl: string): Promise<string | null> => {
    const refreshToken = getRefreshToken();

    if (!refreshToken) {
        return null;
    }

    const response = await fetch(`${apiBaseUrl}/admin/refresh`, {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
        return null;
    }

    const payload: TokenPairPayload = await response.json();
    const tokens = extractTokenPair(payload);

    if (!tokens) {
        return null;
    }

    setAuthTokens(tokens.token, tokens.refreshToken);
    return tokens.token;
};

export const refreshAdminAccessToken = async (apiBaseUrl: string): Promise<string | null> => {
    if (!refreshRequestInFlight) {
        refreshRequestInFlight = fetchRefreshedAccessToken(apiBaseUrl)
            .catch(() => null)
            .finally(() => {
                refreshRequestInFlight = null;
            });
    }

    const nextAccessToken = await refreshRequestInFlight;
    if (!nextAccessToken) {
        clearAuthTokens();
    }

    return nextAccessToken;
};
