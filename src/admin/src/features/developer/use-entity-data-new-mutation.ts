import { useFullEntity } from "@/hooks/use-entities";
import { LGQuery } from "@/lib/apollo/builders/LabraGqlApiBuilder/LGQuery";
import { useMemo } from "react";

export const useEntityDataNewMutation = (entityName: string) => {

    // const fullEntity = useFullEntity({ entityName });

    // if (!fullEntity) {
    //     return null;
    // }

    var query = useMemo(() => {
        return getEntityDataNewMutationQuery(entityName);
    }, [entityName]);

    return query;
}

export const getEntityDataNewMutationQuery = (entityName: string) => {

	return LGQuery.create(entityName, {})
				  .select('id')
				  .build();

    // return new GqlDataCREATEMutationBuilder()
    //     .addEntityName(entityName)
    //     .addField('id')
    //     .build({

    //     });
}