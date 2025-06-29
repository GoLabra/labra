import { useCallback, useMemo } from "react";
import { ChangedFullEntity, ChangedNameCaptionEntity, DesignerEdge, DesignerField, NameCaptionEntity } from "@/types/entity";
import { gql, useApolloClient } from "@apollo/client";
import { useEntitiesDesignerChanges } from "./use-designer-entities-changes";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { ADMIN_CONTEXT } from "@/lib/apollo/apolloWrapper";

export const useEntitiesDesigner = () => {

    const client = useApolloClient();
    const graphEntities = useEntities();
    const designerChanges = useEntitiesDesignerChanges();


    const allNameCaptionEntities: ChangedNameCaptionEntity[] = useMemo(() => {
        let allEntities: ChangedNameCaptionEntity[] = graphEntities.entities.map((i) => ({
                name: i.name,
                caption: i.caption,
                designerStatus: 'unchanged'
            }));

        designerChanges.changedEntityNamesAndCaptions.forEach((entity) => {
            const exists = allEntities.find(i => i.name == entity.name);
            if(exists){
                exists.caption = entity.caption ?? exists.caption;
                exists.designerStatus = entity.designerStatus;
                return;
            }
            allEntities.push({...entity})
        });

        return allEntities.sort((a, b) => a.caption!.localeCompare(b.caption!));
    }, [ graphEntities.entities, designerChanges.changedEntityNamesAndCaptions]);


    const allFullEntities: ChangedFullEntity[] = useMemo(() => {
        return allNameCaptionEntities.map((entity) => {

            const allChildrenNames = [... new Set([
                ...graphEntities.fullEntitiesMap[entity.name]?.fields?.map(c => c.name) ?? [], 
                ...graphEntities.fullEntitiesMap[entity.name]?.edges?.map(c => c.name) ?? [], 
                ...designerChanges.changedFullEntitiesMap[entity.name]?.children?.map(c => c.name) ?? []
            ])];

            const children: (DesignerField | DesignerEdge)[] = allChildrenNames.map(childName => {

                const changedChild = designerChanges.changedFullEntitiesMap[entity.name]?.children.find(j => j.name == childName);
                
                if(changedChild){
                    return changedChild;
                }
                const originalField = graphEntities.fullEntitiesMap[entity.name]?.fields.find(j => j.name == childName);

                if(originalField){
                    return {
                        ...originalField!,
                        designerStatus: 'unchanged'
                    };
                }

                const originalEdge = graphEntities.fullEntitiesMap[entity.name]?.edges.find(j => j.name == childName);
                return {
                    ...originalEdge!,
                    designerStatus: 'unchanged'
                };
                
            });

            const displayFieldCaption = designerChanges.changedFullEntitiesMap[entity.name]?.displayFieldCaption ?? graphEntities.fullEntitiesMap[entity.name]?.displayField?.caption ?? 'Id';

            const result: ChangedFullEntity = {
                name: entity.name,
                caption: designerChanges.changedFullEntitiesMap[entity.name]?.caption ?? entity.caption,
                displayFieldCaption: children.filter(i => i.designerStatus != 'deleted')
                                            .find(i => i.caption == displayFieldCaption)?.caption
                                            ?? 'Id',
                children: children,
                fields: children.filter(c => c.__typename == 'Field') as DesignerField[],
                edges: children.filter(c => c.__typename == 'Edge') as DesignerEdge[],
                designerStatus: entity.designerStatus,
                loading: graphEntities.fullEntitiesMap[entity.name]?.loading ?? false
            }
            return result;
        })
    }, [ allNameCaptionEntities, graphEntities.fullEntitiesMap, designerChanges.changedFullEntitiesMap]);

    const allFullEntitiesMap = useMemo(() => {
        return allFullEntities.reduce((acc: Record<string, ChangedFullEntity>, i: ChangedFullEntity) => {
            acc[i.name] = {... i};
            return acc;
        }, {});
    }, [allFullEntities]);

    const updateChild = useCallback((entityName: string, fieldName: string, field: DesignerField | DesignerEdge) => {
        const currentChild = allFullEntitiesMap[entityName].children.find(i => i.name == fieldName);
        
        if(JSON.stringify(currentChild) == JSON.stringify(field)){
            return;
        }

        const addedChild = designerChanges.updateChild(entityName, fieldName, field);

        if(currentChild?.caption == allFullEntitiesMap[entityName].displayFieldCaption){
            // don't set if caption is the same
            if(currentChild?.caption == addedChild.caption){
                return;
            }
            designerChanges.setDisplayField(entityName, addedChild.caption);
        }
    }, [designerChanges.updateChild, allFullEntitiesMap, designerChanges.setDisplayField]);
    
    const deleteEntity = useCallback((entityName: string) => {
        var entity = allNameCaptionEntities.find(i => i.name == entityName);
        if(!entity){
            return;
        }

        designerChanges.deleteEntity(entity);
    }, [graphEntities, allFullEntitiesMap, designerChanges.deleteEntity]);

    const deleteChild = useCallback((entityName: string, fieldName: string) => {
        const field = allFullEntitiesMap[entityName].children.find(i => i.name == fieldName) ?? null;

        if(!field){
            return;
        }

        designerChanges.deleteChild(entityName, fieldName, field);

        const currentChild = allFullEntitiesMap[entityName].children.find(i => i.name == fieldName);
        if(currentChild?.caption == allFullEntitiesMap[entityName].displayFieldCaption){
            designerChanges.setDisplayField(entityName, undefined);
        }

    }, [allFullEntitiesMap, designerChanges.deleteChild]);

    const saveChanges = useCallback((fullEntity: ChangedFullEntity) => {
        const query = designerChanges.getSaveQuery(fullEntity);
        if(!query){
            return;
        }

        client.mutate({ mutation: gql(query.query), variables: query.variables, fetchPolicy: "network-only", context: ADMIN_CONTEXT })
            .then((response) => {
                // clear designer changes
                designerChanges.revertEntityChange(fullEntity.name);
            }).finally(() => {

            });

    }, [allFullEntities, designerChanges.getSaveQuery, designerChanges.revertAllChanges]);

    const setDisplayField = useCallback((entityName: string, displayFieldCaption?: string) => {
        if(graphEntities.fullEntitiesMap[entityName]?.displayField?.caption == displayFieldCaption){
            designerChanges.setDisplayField(entityName, undefined);
            return;
        }
        designerChanges.setDisplayField(entityName, displayFieldCaption);
    }, [designerChanges.setDisplayField, graphEntities.fullEntitiesMap]);

    return useMemo(() => ({
        allEntities: allFullEntities,
        allFullEntitiesMap: allFullEntitiesMap,

        allNameCaptionEntities: allNameCaptionEntities,

        addEntity: designerChanges.addEntity,
        editEntity: designerChanges.editEntity,

        setDisplayField: setDisplayField,

        addChild: designerChanges.addChild,
        updateChild: updateChild,
        
        deleteEntity,
        deleteChild,

        revertAllChanges: designerChanges.revertAllChanges,
        revertEntityChange: designerChanges.revertEntityChange,
        revertChildChange: designerChanges.revertChildChange,

        hasChanges: designerChanges.hasChanges,

        saveChanges
    }), [allFullEntities, allFullEntitiesMap, allNameCaptionEntities, designerChanges.addEntity, designerChanges.editEntity, setDisplayField, designerChanges.addChild, updateChild, deleteEntity, deleteChild, designerChanges.revertAllChanges, designerChanges.revertEntityChange, designerChanges.revertChildChange, designerChanges.hasChanges, saveChanges]);
};


export const useEntityDesignerForEntity = (entityName: string) => {

    const originalFullEntity = useFullEntity({ entityName });
    
    const designerEntity = useEntitiesDesigner();
    const designerChanges = useEntitiesDesignerChanges();

    const nameCaptionEntity = useMemo(() => {
        return designerEntity.allNameCaptionEntities.find(i => i.name == entityName);
    }, [designerEntity.allNameCaptionEntities, entityName]);

    const fullDesignerEntity = useMemo(() => {
        const fullEntity = designerEntity.allFullEntitiesMap[entityName];
        return fullEntity
    }, [designerEntity.allFullEntitiesMap, entityName]);

    const entityStatus = useMemo(() => {
        return designerChanges.changedFullEntitiesMap[entityName]?.designerStatus ?? 'unchanged';
    }, [designerChanges.changedFullEntitiesMap, entityName]);

    const justChanges = useMemo(() => {
        return designerChanges.changedFullEntitiesMap[entityName];
    }, [designerChanges.changedFullEntitiesMap]);
    
    const hasChanges = useMemo(() => {
        return !!justChanges;
    }, [justChanges]);

    const editEntity = useCallback((newEntity: ChangedNameCaptionEntity) => {
        designerEntity.editEntity(entityName, newEntity);
    },[ designerEntity.editEntity, entityName]);

    const deleteEntity = useCallback(() => {
        designerEntity.deleteEntity(entityName);
    }, [designerEntity.deleteEntity, entityName]);

    const addChild = useCallback((child: DesignerField | DesignerEdge) => {
        designerEntity.addChild(entityName, child);
    }, [designerEntity.addChild, entityName]);
    
    const updateChild = useCallback((fieldName: string, child: DesignerField | DesignerEdge) => {
        designerEntity.updateChild(entityName, fieldName, child);
    }, [designerEntity.updateChild, entityName]);
    
    const deleteChild = useCallback((childName: string) => {
        designerEntity.deleteChild(entityName, childName);
    }, [designerEntity.deleteChild, entityName]);

    const setDisplayField = useCallback((displayFieldCaption: string) => {
        designerEntity.setDisplayField(entityName, displayFieldCaption);
    }, [designerEntity.setDisplayField, entityName]);

    const revertChildChange = useCallback((child: DesignerField | DesignerEdge) => {
        designerEntity.revertChildChange(entityName, child.name);
    }, [designerEntity.revertChildChange, entityName]);

    const revertEntityChange = useCallback(() => {
        designerEntity.revertEntityChange(entityName);
    }, [designerEntity.revertEntityChange, entityName]);

    const saveChanges = useCallback(() => {
        designerEntity.saveChanges(fullDesignerEntity);
    }, [designerEntity.saveChanges]);

    return useMemo(() => ({
        originalFullEntity,
        nameCaptionEntity,
        fullDesignerEntity,

        justChanges,
        hasChanges,
        entityStatus,

        editEntity,
        deleteEntity,

        addChild,
        updateChild,
        deleteChild,
        setDisplayField,
        revertChildChange,
        revertEntityChange,

        saveChanges
    }), [ originalFullEntity, nameCaptionEntity, fullDesignerEntity, justChanges, hasChanges, entityStatus, editEntity, deleteEntity, addChild, updateChild, deleteChild, setDisplayField, revertChildChange, revertEntityChange, saveChanges]);
}