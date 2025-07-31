import { useFullEntity } from "@/hooks/use-entities";
import { ApiType } from "@/lib/apollo/apolloWrapper";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { EntityOwner } from "@/lib/apollo/graphql.entities";
import { pascalCase } from "change-case";
import { useMemo } from "react";


export const useEntityDataDeleteMutation = (entityName: string) => {

	const fullEntity = useFullEntity({ entityName });

    var query = useMemo(() => {
        return getEntityDataDeleteMutationQuery(entityName);
    }, [entityName]);

    return useMemo(() => ({
		...query,
		lgQuery: getEntityDataDeleteMutationLGQuery(entityName),
		apiType: fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user' as ApiType
	}), [query, entityName, fullEntity?.owner]);
}

export const getEntityDataDeleteMutationQuery = (entityName: string) => {

	return LGQuery.deleteFrom(entityName)
		.where(GplFilter.field('id', '', '123'))
		.select('id')
		.build();
}

export const getEntityDataDeleteMutationLGQuery = (entityName: string) => {

	return `LGQuery.deleteFrom<${pascalCase(entityName)}>('${entityName}')
	.where(GplFilter.field('id', '', '123'))
	.select('id')`;
}