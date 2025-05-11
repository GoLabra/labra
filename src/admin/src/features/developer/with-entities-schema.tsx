import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntitiesSchema } from "./use-entities-schema";

export const WithEntitiesSchema = (Component: ComponentType<ShowGraphQlQueryProps>) => {
    const WrappedComponent: FC = ( ) => {
        const query = useEntitiesSchema();

        if (!query) {
            return null;
        }

        return (
            <Component
                title="Entities Schema"
                query={query.query}
                variables={query.variables}
                context={query.context}
            />
        );
    };

    WrappedComponent.displayName = `WithEntitySchema(${Component.displayName || Component.name || 'Component'})`;
    return WrappedComponent;
};
