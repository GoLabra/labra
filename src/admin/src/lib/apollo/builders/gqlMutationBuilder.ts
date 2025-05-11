import * as gqlBuilder from 'gql-query-builder'
import { pascalCase } from "change-case";
import { OperationVariables } from "@apollo/client";
import pluralize from "pluralize";

class GqlDataBaseMutationBuilder<TData> {

    protected _entityName: string | null = null;

    private _fields: Array<string> = [];

    public addEntityName = (entityName: string) => {
        this._entityName = entityName;
        return this;
    }

    public addField = (fieldName: string) => {
        this._fields?.push(fieldName);
        return this;
    }

    public addFields = (fieldNames: string[]) => {
        fieldNames.forEach(element => this.addField(element));
        return this;   
    }

    protected buildData = (data: TData): any => {
        throw new Error("Method not implemented.");
    }

    public getOperationName = (): string | null => {
        throw new Error("Method not implemented.");
    }

    public build = (data: TData): { query: string, variables: OperationVariables } | null => {

        const op = this.getOperationName();
        if (!op) {
            return null;
        }

        const dataVariables = this.buildData(data);

        const queryResult = gqlBuilder.mutation({
            operation: op,
            variables: {
                ...dataVariables
            },
            fields: this._fields
        });
        return queryResult;
    }
}

export class GqlDataCREATEMutationBuilder extends GqlDataBaseMutationBuilder<Record<string, any>> {

    public getOperationName = (): string | null => {

        if (!this._entityName) {
            return null;
        }

        return `create${pascalCase(this._entityName)}`;
    }

    protected buildData = (data: Record<string, any>): any => {

        if (!this._entityName) {
            return {};
        }

        const dataType = `Create${pascalCase(this._entityName!)}Input`;

        return {
            data: {
                type: dataType,
                required: true,
                value: data
            }
        }
    }
}

export class GqlDataUPDATEMutationBuilder extends GqlDataBaseMutationBuilder<{ id: string, data: Record<string, any> }> {

    public getOperationName = (): string | null => {

        if (!this._entityName) {
            return null;
        }

        return `update${pascalCase(this._entityName)}`;
    }

    protected buildData = ({ id, data }: { id: string, data: Record<string, any> }): any => {

        if (!this._entityName) {
            return {};
        }

        //cleanup data
        delete data.id;
        delete data.createdBy;
        delete data.updatedBy;
        delete data.createdAt;
        delete data.updatedAt;


        const whereType = `${pascalCase(this._entityName!)}WhereUniqueInput`;
        const dataType = `Update${pascalCase(this._entityName!)}Input`;

        return {
            where: {
                type: whereType,
                required: true,
                value: { id }
            },
            data: {
                type: dataType,
                required: true,
                value: data
            }
        }
    }
}

export class GqlDataDELETEMutationBuilder extends GqlDataBaseMutationBuilder<string> {

    public getOperationName = (): string | null => {

        if (!this._entityName) {
            return null;
        }

        return `delete${pascalCase(this._entityName)}`;
    }

    protected buildData = (id: string): any => {
        if (!this._entityName) {
            return {};
        }

        const dataType = `${pascalCase(this._entityName!)}WhereUniqueInput`;

        return {
            where: {
                type: dataType,
                required: true,
                value: { id: id }
            }
        }
    }
}



export class GqlDataDELETEBulkMutationBuilder extends GqlDataBaseMutationBuilder<Array<string>> {

    public getOperationName = (): string | null => {

        if (!this._entityName) {
            return null;
        }

        return pluralize(`deleteMany${pascalCase(this._entityName)}`);
    }

    protected buildData = (ids: Array<string>): any => {
        if (!this._entityName) {
            return {};
        }

        if (!ids) {
            return {};
        }

        if (!ids.length) {
            return {};
        }

        const dataType = `${pascalCase(this._entityName!)}WhereInput`;

        return {
            where: {
                type: dataType,
                required: true,
                value: {
                    idIn: ids
                }
            }
        }
    }
}
