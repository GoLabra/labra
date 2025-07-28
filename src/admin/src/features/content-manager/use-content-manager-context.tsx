import { createContext, PropsWithChildren, useContext, useMemo } from "react";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { useContentManagerSearchFromQuery } from "@/hooks/use-content-manager-search-from-query";
import { FullEntity } from "@/types/entity";
import { useContentManagerStore } from "@/hooks/use-content-manager-store";

const ContentManagerContext = createContext<{
    entityName: string;
    fullEntity?: FullEntity;
    contentManagerSearch: ReturnType<typeof useContentManagerSearchFromQuery>;
    contentManagerStore: ReturnType<typeof useContentManagerStore>;

    displayFieldName: string;
}>({
    entityName: '',
    fullEntity: null!,
    contentManagerSearch: null!,
    contentManagerStore: null!,

    displayFieldName: null!
});


export const useContentManagerContext = () => {
    return useContext(ContentManagerContext);
};

interface ContextManagerProviderProps {
    entityName: string
}
export const ContentManagerProvider = ({ entityName, children }: PropsWithChildren<ContextManagerProviderProps>) => {

    const fullEntity = useFullEntity({ entityName });
    const contentManagerSearch = useContentManagerSearchFromQuery();

    const contentManagerStore = useContentManagerStore({
		fullEntity: fullEntity!,
		searchState: contentManagerSearch.state,
    });

    const value = useMemo(() => ({
        entityName,
        addNewEntry: () => { },
        fullEntity: fullEntity ?? undefined,
        contentManagerSearch: contentManagerSearch,
        contentManagerStore: contentManagerStore,

        displayFieldName: fullEntity?.displayField?.name ?? 'id'
    }), [entityName, contentManagerSearch, contentManagerStore, fullEntity?.displayField?.name]);

    return (
        <ContentManagerContext.Provider value={value}>
            {children}
        </ContentManagerContext.Provider>
    );
};