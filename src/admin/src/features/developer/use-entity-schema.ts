import { useFullEntity } from "@/hooks/use-entities";
import { ENTITY_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { GqlDataDELETEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { useMemo } from "react";


export const entitySchemaQuery = `  
query EntitySchema($name: String!) {
  entity(where: {name: $name}) {
    name
    caption
    owner
    displayField {
      ...FieldProperties
      __typename
    }
    fields {
      ...FieldProperties
      __typename
    }
    edges {
      name
      caption
      required
      relationType
      private
      relatedEntity {
        name
        caption
        __typename
      }
      belongsToCaption
      __typename
    }
    __typename
  }
}
fragment FieldProperties on Field {
  name
  caption
  type
  required
  unique
  defaultValue
  min
  max
  private
  acceptedValues
  __typename
}
`;

export const useEntitySchema = (entityName: string | null) => {

    const variables = useMemo(() => {
        if (!entityName) {
            return null;
        }

        return {
            name: entityName
        }
    }, [entityName]);

    return useMemo(() => ({
        query: entitySchemaQuery,
        variables,
        context: ENTITY_CONTEXT
    }), []);
}
