export const getCookie = (name: string): string => {
    if (typeof document === 'undefined') {
        return '';
    }

    const cookies = document.cookie ? document.cookie.split('; ') : [];
    for (const c of cookies) {
        const [k, ...v] = c.split('=');
        if (k === name) {
            return decodeURIComponent(v.join('='));
        }
    }

    return '';
};
