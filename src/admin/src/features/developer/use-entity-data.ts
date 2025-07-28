import { useFullEntity } from "@/hooks/use-entities";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
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

	let query =  LGQuery.from<any>(fullEntity.name)
						.skip(0)
						.first(10)
						.orderByAscending(fullEntity.displayField!.name)
						.select(...fullEntity.fields.map(i => i.name));
	// select edges
	query = fullEntity.edges.reduce((query: LGQuery<any>, edge) => {
		return query.include(edge.name, q => q.select('id', edge.relatedEntity.displayField.name));
	}, query)
	
	return query.build();


    // return new GqlDataQueryBuilder()
    //     .addEntityName(fullEntity.name)
    //     .setPagination(1, 10)
    //     .setOrder(fullEntity.displayField!.name, 'asc')
    //     .addFields(fullEntity.fields.map(i => i.name))
    //     .addEdges(fullEntity.edges.map(i => ({
    //         name: i.name,
    //         fields: ['id', i.relatedEntity.displayField.name],
    //     })))
    //     //.addAdvancedFilters(searchState.advancedFilters)
    //     // .addAdvancedFilters(getAdvancedFiltersFromGridFilter(searchState.filter, fields))
    //     // .addOrAdvancedFilters(getAdvancedFiltersFromQuery(searchState.query, fields))
    //     .build();
}