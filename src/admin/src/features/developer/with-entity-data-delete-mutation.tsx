import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityDataDeleteMutation } from "./use-entity-data-delete-mutation";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

export const WithEntityDataDeleteMutation = (Component: ComponentType<ShowGraphQlQueryProps>) => {
    const WrappedComponent: FC<{ entityName: string }> = ({ entityName }) => {
        const query = useEntityDataDeleteMutation(entityName);

        if (!query) {
            return null;
        }

        return (
            <Component
                title="Delete Entry"
                query={query.query}
                variables={query.variables}
				lqQuery={query.lgQuery}
				apiType={query.apiType}
            />
        );
    };

    WrappedComponent.displayName = `WithEntityDataDeleteMutation(${Component.displayName || Component.name || 'Component'})`;
    
    return WrappedComponent;
};
