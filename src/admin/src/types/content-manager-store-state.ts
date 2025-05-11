import { Edge, Field } from "@/lib/apollo/graphql.entities";


export interface ContentManagerStoreState {

    data: any
    // schemaLoading: boolean;
    dataLoading: boolean;    
    // entityFields: Field[];
    // entityEdges: Edge[];

    pagesCount: number;
    totalItems: number;
}

export interface UsersStore {
    state: ContentManagerStoreState;

    refresh: () => void;
    addItem: (data: any) => void;
    addItems: (data: any[]) => void;
    updateItem: (id: string, data: any) => void;
    deleteItem: (id: string) => void;
    deleteBulk: (ids: Array<string>) => void;
}
