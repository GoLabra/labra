import { useMemo, useState, useCallback } from 'react';

export function useTabs(defaultValue: string) {
    const [value, setValue] = useState<string>(defaultValue);

    const onChange = useCallback((event: any, newValue: string) => {
        setValue(newValue);
    }, []);

    return useMemo(() => ({
        value,
        setValue,
        onChange
    }), [onChange, value]);
}
