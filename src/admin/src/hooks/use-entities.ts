import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";
import { Edge, Entity, Field } from "@/lib/apollo/graphql.entities";
import { FullEntity, NameCaptionEntity } from "@/types/entity";
import { gql, makeVar, useLazyQuery, useQuery, useReactiveVar } from "@apollo/client";
import { useCallback, useEffect, useMemo, useState } from "react";

export const GetEntityFirstLevelSchema = gql`
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
     }

    query GetEntityFirstLevelSchema($name: String!) {
        entity(where: {name: $name}) {
            name
            caption
            owner
            displayField {
               ...FieldProperties
            }
            fields {
               ...FieldProperties
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
                    displayField {
                        name
                    }
                },
                belongsToCaption
            }
        }
    }`

export const GetEntitiesNameCaption = gql`query GetEntitiesNameCaption {
    entities{
      name
      caption
  }
}`


export const nameCaptionEntitiesVar = makeVar<NameCaptionEntity[]>([]);
export const fullEntitiesMapVar = makeVar<Record<string, FullEntity>>({});


export const useEntities = () => {
    const nameCaptionEntities = useReactiveVar<NameCaptionEntity[]>(nameCaptionEntitiesVar);
    const fullEntitiesMap = useReactiveVar<Record<string, FullEntity>>(fullEntitiesMapVar);

    const { data, loading: loadingNameCaptionEntities } = useQuery(GetEntitiesNameCaption, {context: ADMIN_CONTEXT}); //useGetEntitiesNameCaptionQuery();

    const [loadFullEntity] = useLazyQuery<{entity: Entity}>(GetEntityFirstLevelSchema, {context: ADMIN_CONTEXT}); //useGetEntityFirstLevelSchemaLazyQuery(); 

    // set entities
    useEffect(() => {
        if (!data) {
            return;
        }

        nameCaptionEntitiesVar((data!.entities ?? []));

    }, [data]);

    const onLoadFullEntity = useCallback((entityName: string) => {

        if (fullEntitiesMap[entityName]) {
            return;
        }

        if (!nameCaptionEntities.find(i => i.name == entityName)) {
            return;
        }

        const emptyEntityLoading: FullEntity = {
            name: entityName,
            caption: null!,
            owner: null!,
            displayField: null!,
            fields: [],
            edges: [],
            loading: true,
        };

        fullEntitiesMapVar({
            ...fullEntitiesMapVar(),
            [entityName]: emptyEntityLoading
        });

        loadFullEntity({ variables: { name: entityName }, fetchPolicy: 'network-only' }).then((apolloResponse) => {

            const { data } = apolloResponse;

            if (!data) {
                return;
            }

            if (!data.entity) {
                return;
            }

            fullEntitiesMapVar({
                ...fullEntitiesMap,
                [data.entity.name]: { 
                    name: data.entity.name,
                    caption: data.entity.caption,
                    owner: data.entity.owner,
                    displayField: data.entity.displayField,
                    fields: data.entity?.fields as Array<Field> ?? [],
                    edges: data.entity?.edges?.filter(i => i.name != "refCreatedBy")
                                            .filter(i => i.name != "refUpdatedBy")
                                            .filter(i => i.caption != "Ref Updated By") as Array<Edge> ?? [],
                    loading: false,
                }
            });
        });

    }, [loadFullEntity, fullEntitiesMap, nameCaptionEntities]);

    return useMemo(() => ({
        entities: nameCaptionEntities,
        fullEntitiesMap: fullEntitiesMap,
        loadingNameCaptionEntities: loadingNameCaptionEntities,

        onLoadFullEntity: onLoadFullEntity,

    }), [nameCaptionEntities, fullEntitiesMap, loadingNameCaptionEntities, onLoadFullEntity]);
};


interface useFullEntityProps {
    entityName: string;
}
export const useFullEntity = (props: useFullEntityProps): FullEntity | null => {
    const entities = useEntities();

    useEffect(() => {
        entities.onLoadFullEntity(props.entityName);
    }, [props.entityName, entities]);

    const fullEntitiesMap = entities.fullEntitiesMap[props.entityName];
    return fullEntitiesMap;
}


interface useFullEntitiesProps {
    lazy?: boolean;
}
export const useFullEntities = (props: useFullEntitiesProps) => {

    const entities = useEntities();

    useEffect(() => {
        if(props.lazy){
            return;
        }
        entities.entities.forEach(i => entities.onLoadFullEntity(i.name));
    }, [entities, props.lazy]);

    const allLoaded = useMemo(() => {
        return Object.keys(entities.fullEntitiesMap).length == entities.entities.length;
    }, [entities]);

    const loadAll = useCallback(() => {
        entities.entities.forEach(i => entities.onLoadFullEntity(i.name));
    }, [entities]);

    return useMemo(() => ({
        entities: entities.fullEntitiesMap,
        firestLevelLoaded: !!entities.entities.length,
        allLoaded: allLoaded,
        loadAll
    }), [entities, allLoaded, loadAll]);
}
