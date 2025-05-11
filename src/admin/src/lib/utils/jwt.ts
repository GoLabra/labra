import { jwtDecode } from "jwt-decode";

export const isJwtValid = (accessToken?: string | null): boolean => {

    if(!accessToken) {
        return false;
    }

    const decoded = jwtDecode(accessToken);
    const currentTime = Date.now() / 1000;
    if (decoded.exp && decoded.exp < currentTime) {
        return false;
    }

    return true;
}


export const getJwtSub = (accessToken?: string | null): string | null | undefined => {

    if(!accessToken) {
        return null;
    }

    const decoded = jwtDecode(accessToken);
    return decoded.sub;
}

export const getJwtRole = (accessToken?: string | null): string | null | undefined => {

    if(!accessToken) {
        return null;
    }

    const decoded = jwtDecode<{role: string}>(accessToken);
    return decoded.role;
}