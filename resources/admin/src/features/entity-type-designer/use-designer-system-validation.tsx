import { useCallback } from "react";
import { useEntityDesignerForEntity } from "./use-designer-entities"; 
import { z } from "zod";
import { useParamSingleValue } from "@/hooks/useParam";
import { ChangedFullEntity } from "@/types/entity";
import { ChildTypeDescriptor } from "./designer-field-map";


export const checkChildCaptionExists = (caption: string, editId: string | undefined, entity: ChangedFullEntity) => {
    const normalizedCaption = caption.toLocaleLowerCase();

    const entityChildren = entity?.children ?? [];
    const existingField = entityChildren?.find(i => i.caption.toLocaleLowerCase() == normalizedCaption && i.name != editId);
    return !existingField;
}

export const useSchemaEffectEntityParams = (editId?: string, schemaEffect?: ChildTypeDescriptor['schemaEffect']) => {

    const entityName = useParamSingleValue('entity-id');
    const entitiesDesigner = useEntityDesignerForEntity(entityName!);

    return useCallback((schema: z.ZodObject<any, any>) => {
        return schemaEffect?.(schema, editId, entitiesDesigner.fullDesignerEntity) ?? schema;
    }, [entityName, editId, schemaEffect, entitiesDesigner.fullDesignerEntity?.children]);
}