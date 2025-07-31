import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityDataUpdateMutation } from "./use-entity-data-update-mutation";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

export const WithEntityDataUpdateMutation = (Component: ComponentType<ShowGraphQlQueryProps>) => {

    const WrappedComponent: FC<{ entityName: string }> = ({ entityName }) => {

        var query = useEntityDataUpdateMutation(entityName);

        if (!query) {
            return null;
        }

        return (
            <Component
                title="Update Entry"
                query={query.query}
                variables={query.variables}
				lqQuery={query.lgQuery}
				apiType={query.apiType}
            />
        )
    };

    WrappedComponent.displayName = `WithEntityData(${Component.displayName || Component.name || 'Component'})`;
    
    return WrappedComponent;
};
