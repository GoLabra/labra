import { ChangedFullEntity, ChangedNameCaptionEntity, DesignerEdge, DesignerField } from "@/types/entity";
import { makeVar } from "@apollo/client";
import { Edge, Field } from "@/lib/apollo/graphql.entities";

export const modifiedEntityVar = makeVar<ChangedNameCaptionEntity[]>([]);
export const modifiedDisplayFieldVar = makeVar<Record<string, string>>({});
export const modifiedEntityChildrenVar = makeVar<Record<string, (DesignerField | DesignerEdge)[]>>({});

export interface IEntityDesignerStrategy {
    getFlowName(): string,
    getEntities(): ChangedNameCaptionEntity[];
    setDisplayField(entityName: string, displayField?: string): void;
    addEntity(newEntity: ChangedNameCaptionEntity): void;
    editEntity(name: string, newEntity: ChangedNameCaptionEntity): void;
    deleteEntity(entity?: ChangedNameCaptionEntity): void;

    addField(entityName: string, newField: DesignerField|DesignerEdge): DesignerField | DesignerEdge;
    updateChild(entityName: string, fieldName: string, updatedField: DesignerField | DesignerEdge): DesignerField | DesignerEdge;
    deleteChild(entityName: string, fieldName: string, deletedField?: Field | Edge | null): void;

    getSaveQuery(fullEntity:ChangedFullEntity): { query: string, variables: Record<string, any> } | null
}