import { ContentManagerStoreState } from "@/types/content-manager-store-state";
import { useCallback, useMemo } from "react";

export const useContentManagerIds = (storeState: ContentManagerStoreState): {storeIds: [], getId: (row: any) => string} => {
    const storeIds = useMemo(() => {
            if (!storeState.data) {
                return [];
            }
            
            return storeState.data.map((data: any) => data.id);
        }, [storeState]);

    const getId = useCallback((row: any) => {
        return row?.id;
    }, []);

    return {
        storeIds,
        getId
    }
};

