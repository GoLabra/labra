import { useFullEntity } from "@/hooks/use-entities";
import { ApiType } from "@/lib/apollo/apolloWrapper";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { EntityOwner } from "@/lib/apollo/graphql.entities";
import { FullEntity } from "@/types/entity";
import { pascalCase } from "change-case";
import { useMemo } from "react";

export const useEntityData = (entityName: string) => {

    const fullEntity = useFullEntity({ entityName });

    var query = useMemo(() => {
        if (!fullEntity || fullEntity.loading) {
            return;
        }

        return getEntityDataQuery(fullEntity);
    }, [fullEntity]);

	var lgQuery = useMemo(() => {
        if (!fullEntity || fullEntity.loading) {
            return;
        }

        return getEntityDataLGQuery(fullEntity);
    }, [fullEntity]);

    return useMemo(() => ({
		...query,
		lgQuery: lgQuery,
		apiType: fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user' as ApiType
	}), [query, entityName, fullEntity?.owner]);
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
}

export const getEntityDataLGQuery = (fullEntity: FullEntity): string => {

	let query = `LGQuery.from<${pascalCase(fullEntity.name)}>('${fullEntity.name}')
	.skip(0)
	.first(10)
	.orderByAscending('${fullEntity.displayField!.name}')
	.select(${fullEntity.fields.map(i => `'${i.name}'`).join(', ')})`;
	// select edges

	for(var edge of fullEntity.edges){
		query += `
	.include('${edge.name}', q => q.select('id', '${edge.relatedEntity.displayField.name}'))`;
	}

	return query;
}