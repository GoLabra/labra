import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityData } from "./use-entity-data";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

export const WithEntityData = (Component: ComponentType<ShowGraphQlQueryProps>) => {

    const WrappedComponent: FC<{ entityName: string }> = ({ entityName }) => {

        var query = useEntityData(entityName);

        if (!query?.query) {
            return null;
        }

        return (
            <Component
                title="Get Entries"
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
