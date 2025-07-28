import { useFullEntity } from "@/hooks/use-entities";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { GplFilter } from "@/lib/apollo/builders/LabraGqlApiBuilder/types/types";
import { useMemo } from "react";

export const useEntityDataUpdateMutation = (entityName: string) => {

    // const fullEntity = useFullEntity({ entityName });

    // if (!fullEntity) {
    //     return null;
    // }

    var query = useMemo(() => {
        return getEntityDataUpdateMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataUpdateMutationQuery = (entityName: string) => {

	return LGQuery.update(entityName, {})
		.where(GplFilter.field('id', '', '123'))
		.select('id')
		.build();
}