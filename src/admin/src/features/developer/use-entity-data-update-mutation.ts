import { useFullEntity } from "@/hooks/use-entities";
import { ApiType } from "@/lib/apollo/apolloWrapper";
import { EntityOwner } from "@/lib/apollo/graphql.entities";
import { pascalCase } from "change-case";
import { GplFilter, LGQuery } from "lg-query";
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
		apiType: fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user' as ApiType
	}), [query, entityName, fullEntity?.owner]);
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