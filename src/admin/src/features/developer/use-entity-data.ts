import { useFullEntity } from "@/hooks/use-entities";
import { GqlDataUPDATEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { GqlDataQueryBuilder } from "@/lib/apollo/builders/gqlQueryBuilder";
import { FullEntity } from "@/types/entity";
import { useMemo } from "react";

export const useEntityData = (entityName: string) => {

    const fullEntity = useFullEntity({ entityName });

    var query = useMemo(() => {
        if (!fullEntity || fullEntity.loading) {
            return;
        }

        return getEntityDataQuery(fullEntity);
    }, [fullEntity]);

    return query;
}

export const getEntityDataQuery = (fullEntity: FullEntity) => {
    return new GqlDataQueryBuilder()
        .addEntityName(fullEntity.name)
        .setPagination(1, 10)
        .setOrder(fullEntity.displayField!.name, 'asc')
        .addFields(fullEntity.fields.map(i => i.name))
        .addEdges(fullEntity.edges.map(i => ({
            name: i.name,
            fields: ['id', i.relatedEntity.displayField.name],
        })))
        //.addAdvancedFilters(searchState.advancedFilters)
        // .addAdvancedFilters(getAdvancedFiltersFromGridFilter(searchState.filter, fields))
        // .addOrAdvancedFilters(getAdvancedFiltersFromQuery(searchState.query, fields))
        .build();
}