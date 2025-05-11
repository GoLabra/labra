
import { Entity } from '@/lib/apollo/graphql.entities';
import { IEntityDesignerStrategy, modifiedDisplayFieldVar, modifiedEntityChildrenVar, modifiedEntityVar } from './IEntityDesignerStrategy';
import { DesignerField, DesignerEdge, ChangedNameCaptionEntity, ChangedFullEntity, DesignerEntityStatus } from '@/types/entity';
import { generateEntityQuery } from '@/lib/apollo/utils/generateEntityQuery';

export class MultipleEntityDesignerStrategyService implements IEntityDesignerStrategy {

    getFlowName() {
        return 'multi';
    }

    getEntities() {
        return modifiedEntityVar();
    }

    setDisplayField(entityName: string, displayField?: string) {
        //this.resetChangesIfOtherEntity(entityName);

        if (!displayField) {
            modifiedDisplayFieldVar(Object.fromEntries(
                Object.entries(modifiedDisplayFieldVar()).filter(([key]) => key !== entityName)
            ));
            return;
        }

        modifiedDisplayFieldVar({
            ...modifiedDisplayFieldVar(),
            [entityName]: displayField
        });
    }

    addEntity(entity: ChangedNameCaptionEntity) {
        const newDesignerEntity: ChangedNameCaptionEntity = {
            ...entity,
            designerStatus: 'new'
        };
        modifiedEntityVar([
            ...modifiedEntityVar(),
            newDesignerEntity]
        );
        // modifiedDisplayFieldVar({});
        // modifiedEntityChildrenVar({});
    }

    editEntity(entityName: string, editedEntity: ChangedNameCaptionEntity) {
        //reset edited entity if other entity
        //this.resetChangesIfOtherEntity(entityName);

        modifiedEntityVar([
            ...modifiedEntityVar().filter(i => i.name != entityName),
            {
                ...editedEntity,
                name: entityName,
                designerStatus: 'edited',
            }]);
    }

    deleteEntity(entity: ChangedNameCaptionEntity) {

        if (entity.designerStatus == 'new') {
            modifiedEntityVar([
                ...modifiedEntityVar().filter(i => i.name != entity.name),
            ]);
            modifiedDisplayFieldVar({
                ...Object.fromEntries(
                    Object.entries(modifiedDisplayFieldVar()).filter(([key]) => key !== entity.name)
                )
            });
            modifiedEntityChildrenVar({
                ...Object.fromEntries(
                    Object.entries(modifiedEntityChildrenVar()).filter(([key]) => key !== entity.name)
                )
            });
            return;
        }

        const newDesignerEntity: ChangedNameCaptionEntity = {
            ...entity!,
            designerStatus: 'deleted',
        };
        modifiedEntityVar([
            ...modifiedEntityVar().filter(i => i.name != newDesignerEntity.name),
            newDesignerEntity]);
        modifiedDisplayFieldVar({
            ...Object.fromEntries(
                Object.entries(modifiedDisplayFieldVar()).filter(([key]) => key !== entity.name)
            )
        });
        modifiedEntityChildrenVar({
            ...Object.fromEntries(
                Object.entries(modifiedEntityChildrenVar()).filter(([key]) => key !== entity.name)
            )
        });
    }

    addField(entityName: string, newField: DesignerField | DesignerEdge): DesignerField | DesignerEdge {

        //this.resetChangesIfOtherEntity(entityName);

        const entitiesFields = modifiedEntityChildrenVar();

        const finalNewField: DesignerField | DesignerEdge = {
            ...newField,
            designerStatus: 'new'
        };

        modifiedEntityChildrenVar({
            ...modifiedEntityChildrenVar(),
            [entityName]: [
                ...entitiesFields[entityName] ?? [],
                finalNewField
            ]
        });

        return finalNewField;
    }

    updateChild(entityName: string, fieldName: string, updatedField: DesignerField | DesignerEdge): DesignerField | DesignerEdge {

        //this.resetChangesIfOtherEntity(entityName);

        const entitiesFields = modifiedEntityChildrenVar();
        const field = entitiesFields[entityName]?.find(i => i.name == fieldName) as DesignerField;

        var newStatus: DesignerEntityStatus = 'edited';
        if (field) {
            newStatus = field.designerStatus;
        }

        const finalUpdatedField: DesignerField | DesignerEdge = {
            ...updatedField,
            name: fieldName,
            designerStatus: newStatus
        }

        modifiedEntityChildrenVar({
            ...modifiedEntityChildrenVar(),
            [entityName]: [
                ...entitiesFields[entityName]?.filter(i => i.name != fieldName) ?? [],
                finalUpdatedField
            ]
        });

        return finalUpdatedField;
    }

    deleteChild(entityName: string, fieldName: string, field: DesignerField | DesignerEdge) {

        //this.resetChangesIfOtherEntity(entityName);


        if (field.designerStatus == 'new') {
            let children = (modifiedEntityChildrenVar()[entityName] ?? []).filter(i => i.name != fieldName);

            modifiedEntityChildrenVar({
                ...(children.length ? {
                    [entityName]: children
                } : {})
            });
            return;
        }

        if (field.designerStatus == 'edited' || field.designerStatus == 'unchanged') {
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

        if (!fullEntity) {
            return null;
        }

        const query = generateEntityQuery(fullEntity);
        return query;
    }

}
