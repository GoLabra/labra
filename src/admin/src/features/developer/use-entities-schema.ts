import { useFullEntity } from "@/hooks/use-entities";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { GqlDataDELETEMutationBuilder } from "@/lib/apollo/builders/gqlMutationBuilder";
import { useMemo } from "react";


export const entitiesSchemaQuery = `
query EntitiesSchema {
  entities {
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

export const useEntitiesSchema = () => {

    return useMemo(() => ({ 
        query: entitiesSchemaQuery,
        variables: null,
        context: ADMIN_CONTEXT
    }), []);
}
