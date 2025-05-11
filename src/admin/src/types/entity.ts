import { Edge, Entity, EntityOwner, Field, Maybe } from "@/lib/apollo/graphql.entities";

export type NameCaptionEntity = {
    name: string;
    caption: string;
}


export type FullEntity = {
    name: string;
    caption?: string;
    owner: EntityOwner;
    displayField?: Field;
    fields: Field[];
    edges: Edge[];
    //children: (Field | Edge)[];

    loading: boolean
}

export type ChangedNameCaptionEntity = {
    //__typename?: 'Entity';
    name: string;
    caption?: string;
    designerStatus: DesignerEntityStatus;
};

export type ChangedFullEntity = {
    name: string;
    caption?: string;
    displayFieldCaption?: string;
    fields: DesignerField[];
    edges: DesignerEdge[];
    children: (DesignerField | DesignerEdge)[];
    designerStatus: DesignerEntityStatus;
    loading: boolean;
}


export type DesignerEntityStatus = 'unchanged' | 'new' | 'edited' | 'deleted';



export type DesignerField = Field & {
    designerStatus: DesignerEntityStatus;
};

export type DesignerEdge = Edge & {
    designerStatus: DesignerEntityStatus;
};