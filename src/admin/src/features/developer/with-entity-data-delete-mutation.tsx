import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityDataDeleteMutation } from "./use-entity-data-delete-mutation";

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
                context={null}
            />
        );
    };

    WrappedComponent.displayName = `WithEntityDataDeleteMutation(${Component.displayName || Component.name || 'Component'})`;
    
    return WrappedComponent;
};
