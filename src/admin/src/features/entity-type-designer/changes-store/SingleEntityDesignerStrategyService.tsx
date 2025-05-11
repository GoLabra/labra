import { Entity } from '@/lib/apollo/graphql.entities';
import { IEntityDesignerStrategy, modifiedDisplayFieldVar, modifiedEntityChildrenVar, modifiedEntityVar } from './IEntityDesignerStrategy';
import { ChangedFullEntity, ChangedNameCaptionEntity, DesignerEdge, DesignerEntityStatus, DesignerField } from '@/types/entity';
import { generateEntityQuery } from '@/lib/apollo/utils/generateEntityQuery';


export class SingleEntityDesignerStrategyService implements IEntityDesignerStrategy {

    getFlowName() {
        return 'single';
    }
    

    getEntities() {
        return modifiedEntityVar();
    }

    setDisplayField(entityName: string, displayField?: string) {
       this.resetChangesIfOtherEntity(entityName);

       if(!displayField){
            modifiedDisplayFieldVar({});
            return;
       } else {
           modifiedDisplayFieldVar({
                [entityName]: displayField
           });
       }
    }

    addEntity(entity: ChangedNameCaptionEntity) {
        const newDesignerEntity: ChangedNameCaptionEntity = {
            ...entity,
            designerStatus: 'new'
        };
        modifiedEntityVar([newDesignerEntity]);
        modifiedDisplayFieldVar({});
        modifiedEntityChildrenVar({});
    }

    editEntity(entityName: string, editedEntity: ChangedNameCaptionEntity) {
        //reset edited entity if other entity
        this.resetChangesIfOtherEntity(entityName);
        
        modifiedEntityVar([{
            ...editedEntity,
            name: entityName
        }]);
    }

    deleteEntity(entity?: ChangedNameCaptionEntity) {

        if(!entity) {
            modifiedEntityVar([]);
            modifiedDisplayFieldVar({});
            modifiedEntityChildrenVar({});
            return;
        }

        const newDesignerEntity: ChangedNameCaptionEntity = {
            ... entity!,
            designerStatus: 'deleted',
        };
        modifiedEntityVar([newDesignerEntity]);
        modifiedDisplayFieldVar({});
        modifiedEntityChildrenVar({});
    }

    addField(entityName: string, newField: DesignerField|DesignerEdge): DesignerField | DesignerEdge {

        this.resetChangesIfOtherEntity(entityName);

        const entitiesFields = modifiedEntityChildrenVar();

        const finalNewField: DesignerField | DesignerEdge  = {
            ...newField,
            designerStatus: 'new'
        };

        modifiedEntityChildrenVar({
            [entityName]: [
                ...entitiesFields[entityName] ?? [],
                finalNewField
            ]
        });

        return finalNewField;
    }

    updateChild(entityName: string, fieldName: string, updatedField: DesignerField | DesignerEdge): DesignerField | DesignerEdge {

        this.resetChangesIfOtherEntity(entityName);
        
        const entitiesFields = modifiedEntityChildrenVar();
        const field = entitiesFields[entityName]?.find(i => i.name == fieldName) as DesignerField;

        var newStatus: DesignerEntityStatus = 'edited';
        if (field) {
            newStatus = field.designerStatus;
        }

        const finalUpdatedField: DesignerField | DesignerEdge  = {
            ...updatedField,
            name: fieldName,
            designerStatus: newStatus
        }

        modifiedEntityChildrenVar({
            [entityName]: [
                ...entitiesFields[entityName]?.filter(i => i.name != fieldName) ?? [],
                finalUpdatedField
            ]
        });

        return finalUpdatedField;
    }

    deleteChild(entityName: string, fieldName: string, field: DesignerField | DesignerEdge) {

        this.resetChangesIfOtherEntity(entityName);
        
        
        if(field.designerStatus == 'new'){
            let children = (modifiedEntityChildrenVar()[entityName] ?? []).filter(i => i.name != fieldName);

            modifiedEntityChildrenVar({
                ...(children.length ? {
                    [entityName]: children
                } : {})
            });
            return;
        }

        if(field.designerStatus == 'edited' || field.designerStatus == 'unchanged'){
            modifiedEntityChildrenVar({
                [entityName]: [
                    ...(modifiedEntityChildrenVar()[entityName] ?? []).filter(i => i.name != fieldName),
                    {
                        ...field,
                        designerStatus: 'deleted'
                    }
                ]
            });
            return;
        }

    }

    getSaveQuery(fullEntity: ChangedFullEntity): { query: string, variables: Record<string, any> } | null {

        if(!fullEntity){
            return null;
        }

        const query = generateEntityQuery(fullEntity);
        return query;
    }


    private resetChangesIfOtherEntity = (entityName: string) => {
        if(!this.needsReset(entityName)){
            return;
        }
        
        modifiedEntityVar([]);
        modifiedDisplayFieldVar({});
        modifiedEntityChildrenVar({});
    }

    private needsReset = (entityName: string) => {
        // check in changed entities
        if(modifiedEntityVar().length){

            if(modifiedEntityVar().length > 1){
                return true;
            }

            if(modifiedEntityVar()[0].name != entityName){
                return true;
            }
        }

        // check in changed display fields
        const changedDisplayField = Object.keys(modifiedDisplayFieldVar())
        if(changedDisplayField.length){
            if(changedDisplayField.length > 1){
                return true;
            }

            if(changedDisplayField[0] != entityName){
                return true;
            }
        }

        //check in changed children
        const changedChildren = Object.keys(modifiedEntityChildrenVar())
        if(changedChildren.length){
            if(changedChildren.length > 1){
                return true;
            }

            if(changedChildren[0] != entityName){
                return true;
            }
        }

        return false;
    }
}
