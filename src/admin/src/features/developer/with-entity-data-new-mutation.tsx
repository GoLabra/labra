import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntityDataNewMutation } from "./use-entity-data-new-mutation";

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
                context={null}
            />
        )
    };

    WrappedComponent.displayName = `WithEntityDataNewMutation(${Component.displayName || Component.name || 'Component'})`;
    
    return WrappedComponent;
};
