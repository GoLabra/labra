import { createContext, PropsWithChildren, useContext, useMemo } from "react";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { useContentManagerSearchFromQuery } from "@/hooks/use-content-manager-search-from-query";
import { useContentManagerStore } from "@/hooks/use-content-manager-store";
import { FullEntity } from "@/types/entity";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { getAdvancedFiltersFromQuery } from "@/lib/utils/get-filters-from-query";

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
        entityName: entityName,

        page: contentManagerSearch.state.page,
        rowsPerPage: contentManagerSearch.state.rowsPerPage,
        sortBy: contentManagerSearch.state.sortBy,
        order: contentManagerSearch.state.order,

        fields: useMemo(() => fullEntity?.fields.map(i => i.name), [fullEntity?.fields]),
        edges: useMemo(() => fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
                                                .filter(i => i.relationType !== 'ManyToOne')
                                                .filter(i => i.relationType !== 'Many').map(i => ({
                                                    name: i.name,
                                                    fields: ['id', i.relatedEntity.displayField.name],
                                                })), [fullEntity?.edges]),

        filters: useMemo(() => getAdvancedFiltersFromGridFilter(contentManagerSearch.state.filter), [contentManagerSearch.state.filter]),
        orFilters: useMemo(() => getAdvancedFiltersFromQuery(contentManagerSearch.state.query, fullEntity?.fields ?? []), [contentManagerSearch.state.query, fullEntity?.fields])
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