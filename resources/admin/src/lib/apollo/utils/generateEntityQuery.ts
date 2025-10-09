import { ChangedFullEntity } from "@/types/entity";
import { GqlEntityMutationBuilder } from "../builders/gqlEntityMutationBuilder";

export const generateEntityQuery = (entity: ChangedFullEntity) => {

    if(entity.designerStatus == 'unchanged'){
        return null;
    }

    var entityBuilder = new GqlEntityMutationBuilder()
                            .addEntityCaption(entity.caption ?? 'Id')
                            .addEntityDisplayField(entity.displayFieldCaption!);

    // fields                            
    const newFields = entity.fields.filter(i => i.designerStatus == 'new') ?? [];
    entityBuilder.addNewFields(newFields);

    const editedFields = entity.fields.filter(i => i.designerStatus == 'edited') ?? [];
    entityBuilder.addEditedFields(editedFields);

    const deletedFields = entity.fields.filter(i => i.designerStatus == 'deleted').map(i => i.name) ?? [];
    entityBuilder.addDeletedFields(deletedFields);

    //edges
    const newEdges = entity.edges.filter(i => i.designerStatus == 'new') ?? [];
    entityBuilder.addNewEdges(newEdges);

    const editedEdges = entity.edges.filter(i => i.designerStatus == 'edited') ?? [];
    entityBuilder.addEditedEdges(editedEdges);

    const deletedEdges = entity.edges.filter(i => i.designerStatus == 'deleted').map(i => i.name) ?? [];
    entityBuilder.addDeletedEdges(deletedEdges);

    switch(entity.designerStatus){
        case "new": return entityBuilder.buildNewQuery();
        case "edited": return entityBuilder.buildUpdateQuery();
        case "deleted": return entityBuilder.buildDeleteQuery();
        default: return null;
    }
}