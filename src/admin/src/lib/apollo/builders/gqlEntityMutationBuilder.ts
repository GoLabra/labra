import { DesignerEdge, DesignerField } from "@/types/entity";


const createNewEntityQuery = `mutation createNewEntityQuery($data:  CreateEntityInput!) {
    createEntity(data: $data) {
        name
        caption
    }
}`;

const updateEntityQuery = `mutation updateEntityQuery($where: EntityWhereUniqueInput!, $data: UpdateEntityInput!) {
    updateEntity(where: $where, data: $data){
        name
        caption
    }
}`;

const deleteEntityQuery = `mutation deleteEntityQuery($where: EntityWhereUniqueInput!) {
    deleteEntity(where: $where){
        name
        caption
    }
}`;

export class GqlEntityMutationBuilder {

    private _entityCaption: string | null = null;
    private _entityDisplayFieldCaption: string | null = null;

    private _newFields: DesignerField[] = [];
    private _editedFields: DesignerField[] = [];
    private _deletedFields: string[] = [];

    private _newEdges: DesignerEdge[] = [];
    private _editedEdges: DesignerEdge[] = [];
    private _deletedEdges: string[] = [];


    //#region mutations
    public addEntityCaption = (entityCaption: string) => {
        this._entityCaption = entityCaption;
        return this;
    }

    public addEntityDisplayField = (entityDisplayFieldCaption: string) => {
        this._entityDisplayFieldCaption = entityDisplayFieldCaption;
        return this;
    }

    public addNewField = (fields: DesignerField) => {
        this._newFields?.push(fields);
        return this;
    }

    public addNewFields = (fields: DesignerField[]) => {
        fields.forEach(element => this.addNewField(element));
        return this;
    }

    public addEditedField = (field: DesignerField) => {
        this._editedFields?.push(field);
        return this;
    }

    public addEditedFields = (fields: DesignerField[]) => {
        fields.forEach(element => this.addEditedField(element));
        return this;
    }

    public addDeletedField = (fieldName: string) => {
        this._deletedFields?.push(fieldName);
        return this;
    }

    public addDeletedFields = (fieldNames: string[]) => {
        fieldNames.forEach(element => this.addDeletedField(element));
        return this;
    }

    public addNewEdge = (edge: DesignerEdge) => {
        this._newEdges?.push(edge);
        return this;
    }

    public addNewEdges = (edges: DesignerEdge[]) => {
        edges.forEach(element => this.addNewEdge(element));
        return this;
    }

    public addEditedEdge = (edge: DesignerEdge) => {
        this._editedEdges?.push(edge);
        return this;
    }

    public addEditedEdges = (edges: DesignerEdge[]) => {
        edges.forEach(element => this.addEditedEdge(element));
        return this;
    }

    public addDeletedEdge = (edgeName: string) => {
        this._deletedEdges?.push(edgeName);
        return this;
    }

    public addDeletedEdges = (edgeNames: string[]) => {
        edgeNames.forEach(element => this.addDeletedEdge(element));
        return this;
    }
    //#endregion mutation


    private buildCreateFields = () => {

        return {
            create:
                this._newFields.map(i => ({
                    caption: i.caption,
                    type: i.type,
                    required: i.required,
                    unique: i.unique,
                    defaultValue: i.defaultValue,
                    min: i.min,
                    max: i.max,
                    private: i.private,
                    acceptedValues: i.acceptedValues
                }))

        };
    }

    private buildEditedFields = () => {

        return {
            update:
                this._editedFields.map(i => (
                    {
                        where: {
                            name: i.name,
                        },
                        data: {
                            caption: i.caption,
                            required: i.required,
                            unique: i.unique,
                            defaultValue: i.defaultValue,
                            min: i.min,
                            max: i.max,
                            private: i.private,
                            acceptedValues: i.acceptedValues
                        }
                    }))
        };
    }

    private buildDeletedFields = () => {

        return {
            delete:
                this._deletedFields.map(i => ({
                    name: i,
                }))
        };
    }

    private buildCreateEdges = () => {

        return {
            create:
                this._newEdges.map(i => ({
                    caption: i.caption,
                    relatedEntity: {
                        connect: {
                            caption: i.relatedEntity.caption
                        }
                    },
                    belongsToCaption: i.belongsToCaption,
                    required: i.required,
                    relationType: i.relationType,
                    private: i.private
                }))
        };
    }

    private buildEditedEdges = () => {
        return {
            update:
                this._editedEdges.map(i => (
                    {
                        where: {
                            name: i.name,
                        },
                        data: {
                            caption: i.caption,
                            required: i.required,
                            relationType: i.relationType,
                            private: i.private
                        }
                    }
                ))

        };
    }

    private buildDeletedEdges = () => {

        return {
            delete:
                this._deletedEdges.map(i => ({
                    name: i,
                }))
        };
    }

    public buildNewQuery = (): { query: string, variables: Record<string, any>} | null => {

        return {
            query: createNewEntityQuery,
            variables: {
                data: {
                    caption: this._entityCaption,
                    displayField: {
                        caption: this._entityDisplayFieldCaption!,
                    },

                    fields: {
                        ...this.buildCreateFields()
                    },
                    edges: {
                        ...this.buildCreateEdges()
                    },
                }
            }
        };
    }

    public buildUpdateQuery = (): { query: string, variables: Record<string, any> } | null => {

        return {
            query: updateEntityQuery,
            variables: {
                where: {
                    caption: this._entityCaption,
                },
                data: {
                    caption: this._entityCaption,
                    displayField: {
                        caption: this._entityDisplayFieldCaption!,
                    },

                    fields: {
                        ...this.buildCreateFields(),
                        ...this.buildEditedFields(),
                        ...this.buildDeletedFields()
                    },
                    edges: {
                        ...this.buildCreateEdges(),
                        ...this.buildEditedEdges(),
                        ...this.buildDeletedEdges()
                    },
                }
            }
        };
    }

    public buildDeleteQuery = (): { query: string, variables: Record<string, any> } | null => {

        return {
            query: deleteEntityQuery,
            variables: {
                where: {
                    caption: this._entityCaption,
                }
            }
        };
    }

}
