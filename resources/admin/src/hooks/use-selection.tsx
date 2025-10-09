import { useCallback, useEffect, useState } from 'react';

export interface Selection<T> {
    handleDeselectAll: () => void;
    handleDeselectOne: (item: T) => void;
    handleSelectAll: () => void;
    handleSelectOne: (item: T) => void;
    selected: T[];
}
export const useSelection = <T,>(items: T[] = []): Selection<T> => {
    const [selected, setSelected] = useState<T[]>([]);

    useEffect(() => {
        setSelected([]);
    }, [items]);

    const handleSelectAll = useCallback((): void => {
        setSelected([...items]);
    }, [items]);

    const handleSelectOne = useCallback((item: T): void => {
        setSelected((prevState) => [...prevState, item]);
    }, []);

    const handleDeselectAll = useCallback(() => {
        setSelected([]);
    }, []);

    const handleDeselectOne = useCallback((item: T): void => {
        setSelected((prevState) => {
            return prevState.filter((_item) => _item !== item);
        });
    }, []);

    return {
        handleDeselectAll,
        handleDeselectOne,
        handleSelectAll,
        handleSelectOne,
        selected
    };
};
