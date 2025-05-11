import Fields from "gql-query-builder/build/Fields";
import * as gqlBuilder from 'gql-query-builder'
import { pascalCase } from "change-case";
import pluralize from "pluralize";
import { Order } from "mosaic-data-table";
import { AdvancedFilter } from "@/core-features/dynamic-filter/filter";

export type EdgeRequest = {
    name: string;
    fields: string[];
    edges?: EdgeRequest[];    
}

export class GqlDataQueryBuilder {

    private _entityName: string | null = null;
    private _fields: Fields = [];
    private _edges: Fields = [];
    
    private _whereAdvancedFilter: AdvancedFilter[] = [];
    private _orWhereAdvancedFilter: AdvancedFilter[] = [];

    private _sort: string| null = null;
    private _order: Order | null = null;
    private _page: number = 0;
    private _rowsPerPage: number = 0;

    //#region mutations
    public addEntityName = (entityName: string) => {
        this._entityName = entityName;
        return this;
    }

    public setPagination = (page: number, rowsPerPage: number) => {
        this._page = page;
        this._rowsPerPage = rowsPerPage;
        return this;
    }

    public setOrder = (sort: string|null, order: Order|null) => {
        this._sort = sort;
        this._order = order;
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

    public addEdge = (edge: EdgeRequest) => {
        this._edges?.push({ [edge.name]: [
            ...edge.fields,
            ...edge.edges?.map((i) => ({
                [i.name]: i.fields
            })) ?? []
        ] });
        return this;
    }

    public addEdges = (edges: EdgeRequest[]) => {
        edges.forEach(element => this.addEdge(element));
        return this;
    }

    public addFilter = (filter: AdvancedFilter) => {
        this._whereAdvancedFilter?.push(filter);
        return this;
    }

    public addAdvancedFilters = (filters: AdvancedFilter[]) => {
        filters.forEach(element => this.addFilter(element));
        return this;
    }

    public addOrAdvancedFilter = (filter: AdvancedFilter) => {
        this._orWhereAdvancedFilter?.push(filter);
        return this;
    }

    public addOrAdvancedFilters = (filters: AdvancedFilter[]) => {
        filters.forEach(element => this.addOrAdvancedFilter(element));
        return this;
    }

    //#endregion mutation

    public getOperationName = (): string | null => {

        if(!this._entityName){
            return null;
        }

        return pluralize(this._entityName);
    }

    public getOperationConnectionName = (): string => {

        const operation = this.getOperationName();

        return `${operation}Connection`;
    }

    public build = (): {query: string, variables: Record<string, any>} | null => {
        
        const op = this.getOperationName();
        if(!op){
            return null;
        }

        const connectionOp = this.getOperationConnectionName();
        
        const whereVariables = this.buildWhere();
        const paginationVariables = this.buildPagination();
        const sortVariables = this.buildSort();

        const queryResult = gqlBuilder.query([{
            operation: op,
            variables: {
                ...whereVariables,
                ...paginationVariables,
                ...sortVariables
            },
            fields: [ ... this._fields,
                       ... this._edges
                      ]
        }, {
            operation: connectionOp,
            variables: {
                ...whereVariables
            },
            fields: [ "totalCount",
                // TODO: https://github.com/ent/ent/issues/4144
                // GraphQL totalCount on edge is incorrect when edges are queried
                // {
                //     operation: "edges",
                //     fields: [
                //         {
                //             operation: "node",
                //             fields: [
                //                 ... this._fields,
                //                 ... this._edges
                //             ],
                //             variables: {}
                //         }       
                //     ],
                //     variables: {}
                // }
             ]
        }]);
        return queryResult;
    }

    private buildWhere = () => {

        if(!this._entityName){
            return {};
        }

        const andFilterJson = this._whereAdvancedFilter.map((filter: AdvancedFilter) => {
            const suffix = filter.operator?.replace('_', ''); // '_' is used for equals, where there is no suffix, and UI controls need a value
            const filterName = `${filter.property}${suffix}`;
            return {
                [filterName]: filter.value
            };
        });

        const orFilterJson = this._orWhereAdvancedFilter.map((filter: AdvancedFilter) => {

            const suffix = filter.operator?.replace('_', ''); // '_' is used for equals, where there is no suffix, and UI controls need a value
            const filterName = `${filter.property}${suffix}`;
            return {
                [filterName]: filter.value
            };
        });

        //ex: MrAndreiWhereInput, where mrAndrei is the entity name
        const whereType = `${pascalCase(this._entityName)}WhereInput`;

        return {
            where: {
                type: whereType,
                value: {
                    and: andFilterJson,
                    or: orFilterJson
                }
            }
        }
    }

    private buildPagination = () => {
        return {
            skip: (this._page - 1) * this._rowsPerPage,
            first: this._rowsPerPage
        }
    }

    private buildSort = () => {
        if(!this._sort){
            return {};
        }

        const sortType = `${pascalCase(this._entityName!)}Order`;

        return {
            orderBy: {
                type: sortType,
                value: {
                    field: this._sort,
                    direction: this._order == 'asc' ? 'ASC' : 'DESC'
                }
            }
        }
    }
}
