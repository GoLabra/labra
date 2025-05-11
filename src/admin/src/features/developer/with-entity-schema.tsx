import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useEntitySchema } from "./use-entity-schema";

export const WithEntitySchema = (Component: ComponentType<ShowGraphQlQueryProps>) => {
    const WrappedComponent: FC<{ entityName?: string | null }> = ({ entityName }) => {
        const query = useEntitySchema(entityName ?? null);

        if (!query) {
            return null;
        }

        return (
            <Component
                title="Entity Schema"
                query={query.query}
                variables={query.variables}
                context={query.context}
            />
        );
    };

    WrappedComponent.displayName = `WithEntitySchema(${Component.displayName || Component.name || 'Component'})`;
    return WrappedComponent;
};
