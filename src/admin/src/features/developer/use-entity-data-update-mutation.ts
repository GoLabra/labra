import { useFullEntity } from "@/hooks/use-entities";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { EntityOwner } from "@/lib/apollo/graphql.entities";
import { pascalCase } from "change-case";
import { useMemo } from "react";

export const useEntityDataUpdateMutation = (entityName: string) => {

    const fullEntity = useFullEntity({ entityName });

    // if (!fullEntity) {
    //     return null;
    // }

    var query = useMemo(() => {
        return getEntityDataUpdateMutationQuery(entityName);
    }, [entityName]);

    return useMemo(() => ({
		...query,
		lgQuery: getEntityDataUpdateMutationLGQuery(entityName),
		apiType:fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user'
	}), [query]);
}

export const getEntityDataUpdateMutationQuery = (entityName: string) => {

	return LGQuery.update(entityName, {})
		.where(GplFilter.field('id', '', '123'))
		.select('id')
		.build();
}

export const getEntityDataUpdateMutationLGQuery = (entityName: string) => {

	return `LGQuery.update<${pascalCase(entityName)}>('${entityName}', {JSON_DATA})
	.where(GplFilter.field('id', '', '123'))
	.select('id')`;
}