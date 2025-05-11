import { useCallback, useMemo, useRef } from "react";
import { Edge, Entity, Field } from "@/lib/apollo/graphql.entities";
import { useReactiveVar } from "@apollo/client";
import { ChangedFullEntity, ChangedNameCaptionEntity, DesignerEdge, DesignerField } from "@/types/entity";
import { createId } from '@paralleldrive/cuid2';
import { modifiedEntityVar, modifiedDisplayFieldVar, modifiedEntityChildrenVar } from "./changes-store/IEntityDesignerStrategy";
import { MultipleEntityDesignerStrategyService } from "./changes-store/MultipleEntityDesignerStrategyService";

export const useEntitiesDesignerChanges = () => {

    const entityChanges: ChangedNameCaptionEntity[] = useReactiveVar(modifiedEntityVar);
    const displayFieldChanges = useReactiveVar(modifiedDisplayFieldVar);
    const childrenChanges = useReactiveVar(modifiedEntityChildrenVar);
    const strategyService = useRef(new MultipleEntityDesignerStrategyService());

    const changedEntityJustNames = useMemo(() => {
        const changedEntityNames = [...new Set([ // set is here to remove duplicates
            ...entityChanges.map(i => i.name),
            ...Object.keys(displayFieldChanges),
            ...Object.keys(childrenChanges)
        ])];
        return changedEntityNames;
    }, [entityChanges, displayFieldChanges, childrenChanges]);

    const changedEntityNamesAndCaptions: ChangedNameCaptionEntity[] = useMemo(() => {
        return changedEntityJustNames.map((i) => {
            const entity = entityChanges.find(j => j.name == i);
            if(entity){
                return entity;
            }            
            return {
                name: i,
                designerStatus: 'edited',
            }   
        });
    }, [ changedEntityJustNames, entityChanges]);
    
    const changedFullEntities:ChangedFullEntity[] = useMemo(() => {
        const result = changedEntityNamesAndCaptions.map((i) => ({
            ...i,
            displayFieldCaption: displayFieldChanges[i.name],
            fields: childrenChanges[i.name]?.filter(i => i.__typename == 'Field') as DesignerField[] ?? [],
            edges: childrenChanges[i.name]?.filter(i => i.__typename == 'Edge') as DesignerEdge[] ?? [],
            children: childrenChanges[i.name] ?? [],
            loading: false
        }));
        return result;
    }, [ changedEntityNamesAndCaptions, displayFieldChanges, childrenChanges]);

    const changedFullEntitiesMap = useMemo(() => {
        return changedFullEntities.reduce((acc: Record<string, ChangedFullEntity>, i: ChangedFullEntity) => {
            acc[i.name] = {... i};
            return acc;
        }, {});
    }, [changedFullEntities]);

    const setDisplayField = useCallback((entityName: string, displayFieldCaption?: string) => {
        strategyService.current.setDisplayField(entityName, displayFieldCaption);
    }, []);

    const addEntity = useCallback((newEntity: ChangedNameCaptionEntity) => {
        newEntity.name = createId();
        strategyService.current.addEntity(newEntity);
    }, []);

    const editEntity = useCallback((name: string, newEntity: ChangedNameCaptionEntity) => {
        strategyService.current.editEntity(name, newEntity);
    }, []);

    const deleteEntity = useCallback((entity: ChangedNameCaptionEntity) => {
        strategyService.current.deleteEntity(entity);
    }, []);

    const deleteChild = useCallback((entityName: string, fieldName: string, child: DesignerField | DesignerEdge) => {
        strategyService.current.deleteChild(entityName, fieldName, child);
    }, []);

    const addChild = useCallback((entityName: string, field: DesignerField | DesignerEdge) => {
        const newField: DesignerField | DesignerEdge = {
            ...field,
            name: createId()
        }
        return strategyService.current.addField(entityName, newField);
    }, []);

    const updateChild = useCallback((entityName: string, fieldName: string, field: DesignerField | DesignerEdge): DesignerField | DesignerEdge => {
        return strategyService.current.updateChild(entityName, fieldName, field);
    }, []);

    const hasChanges = useMemo((): boolean => {
        return !!changedEntityJustNames.length;
    }, [changedEntityJustNames]);

    const getSaveQuery = useCallback((fullEntity:ChangedFullEntity) => {
        const query = strategyService.current.getSaveQuery(fullEntity);
        return query;
    }, [changedEntityJustNames]);

    const revertEntityChange = useCallback((entityName: string) => {

        // Remove Entity
        modifiedEntityVar(modifiedEntityVar().filter(i => i.name != entityName));

        // Remove display field
        modifiedDisplayFieldVar(Object.fromEntries(
            Object.entries(modifiedDisplayFieldVar()).filter(([key]) => key !== entityName)
        ));

        //Remove all Children for entity
        modifiedEntityChildrenVar(Object.fromEntries(
            Object.entries(modifiedEntityChildrenVar()).filter(([key]) => key !== entityName)
        ));

    }, [modifiedEntityVar(), modifiedEntityChildrenVar()]);

    const revertChildChange = useCallback((entityName: string, fieldName: string) => {

        if(!fieldName){
            return;
        }

        let allChildren = modifiedEntityChildrenVar()[entityName]?.filter(i => i.name != fieldName) ?? [];

        if(!allChildren.length){
            modifiedEntityChildrenVar(Object.fromEntries(
                Object.entries(modifiedEntityChildrenVar()).filter(([key]) => key !== entityName)
            ));
            return;
        }

        modifiedEntityChildrenVar({
            ...modifiedEntityChildrenVar(),
            [entityName]: allChildren
        });

        // // update display field
        // const entities = entityService.current.getEntities(); 
        // let entity = entities.length ? entities[0] : null;
        // if(entity){
        //     if(entity.displayField?.name == fieldName){
        //         entityService.current.setDisplayField(entityName, undefined);
        //     }
        // }
        // // update display field
        // const entities = entityService.current.getEntities(); 
        // let entity = entities.length ? entities[0] : null;
        // if(entity){
        //     if(entity.displayField?.name == fieldName){
        //         entityService.current.setDisplayField(entityName, undefined);
        //     }
        // }

    }, [modifiedEntityChildrenVar()]);

    const revertAllChanges = useCallback(() => {
        modifiedEntityVar([]);
        modifiedDisplayFieldVar({});
        modifiedEntityChildrenVar({});
    }, [modifiedEntityVar(), modifiedEntityChildrenVar()]);

    return {
        changedEntityJustNames,
        changedEntityNamesAndCaptions,
        changedFullEntities,
        changedFullEntitiesMap,
        // changedEntityNames,
        // changedEntities,
        // changedDisplayFields,

        childrenChanges,

        setDisplayField,

        addEntity,
        editEntity,
        addChild,

        updateChild,

        deleteEntity,
        deleteChild,

        revertAllChanges,
        revertEntityChange,
        revertChildChange,

        hasChanges,

        getSaveQuery
    };
};
