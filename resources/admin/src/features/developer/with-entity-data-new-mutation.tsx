import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityDataNewMutation } from "./use-entity-data-new-mutation";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

export const WithEntityDataNewMutation = (Component: ComponentType<ShowGraphQlQueryProps>) => {

    const WrappedComponent: FC<{ entityName: string }> = ({ entityName }) => {

        var query = useEntityDataNewMutation(entityName);
                    
        if (!query) {
            return null;
        }

        return (
            <Component
                title="Create new Entry"
                query={query.query}
                variables={query.variables}
				lqQuery={query.lgQuery}
				apiType={query.apiType}
            />
        )
    };

    WrappedComponent.displayName = `WithEntityDataNewMutation(${Component.displayName || Component.name || 'Component'})`;
    
    return WrappedComponent;
};
