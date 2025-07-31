import { useFullEntity } from "@/hooks/use-entities";
import { ApiType } from "@/lib/apollo/apolloWrapper";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { EntityOwner } from "@/lib/apollo/graphql.entities";
import { pascalCase } from "change-case";
import { useMemo } from "react";

export const useEntityDataNewMutation = (entityName: string) => {

    const fullEntity = useFullEntity({ entityName });

    var query = useMemo(() => {
        return getEntityDataNewMutationQuery(entityName);
    }, [entityName]);

    return useMemo(() => ({
		...query,
		lgQuery: getEntityDataNewMutationLGQuery(entityName),
		apiType: fullEntity?.owner == EntityOwner.Admin ? 'admin' : 'user' as ApiType
	}), [query, entityName, fullEntity?.owner]);
}

export const getEntityDataNewMutationQuery = (entityName: string) => {

	return LGQuery.create(entityName, {})
				  .select('id')
				  .build();
}

export const getEntityDataNewMutationLGQuery = (entityName: string) => {

	return `LGQuery.create<${pascalCase(entityName)}>('${entityName}', {JSON_DATA})
	.select('id')`;
}