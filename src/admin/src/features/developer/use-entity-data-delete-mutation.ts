import { useFullEntity } from "@/hooks/use-entities";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { useMemo } from "react";


export const useEntityDataDeleteMutation = (entityName: string) => {

    var query = useMemo(() => {
        return getEntityDataDeleteMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataDeleteMutationQuery = (entityName: string) => {

	return LGQuery.deleteFrom(entityName)
		.where(GplFilter.field('id', '', '123'))
		.select('id')
		.build();
}