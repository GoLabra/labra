import { Entity } from "@/lib/apollo/graphql.entities";


export type EntityChanges = {
    add: Entity[];
    update: Entity[];
    delete: Entity[];
}

export const getEmptyEntityChanges = () => {
    return {
        add: [],
        update: [],
        delete: [],
    };
}