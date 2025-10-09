'use client'

import { createContext, PropsWithChildren, use, useContext } from "react";
import { useParamSingleValue } from "./useParam";

export const CurrentEntityName = createContext<string | undefined>(undefined);

export const useCurrentEntityNameContext = () => useContext(CurrentEntityName);

interface CurrentEntityProviderProps {
}
export const CurrentEntityProvider = (props: PropsWithChildren<CurrentEntityProviderProps>) => {

    const entityId = useParamSingleValue('entity-id', '');

    return (
        <CurrentEntityName.Provider value={entityId}>
            {props.children}
        </CurrentEntityName.Provider>
    )
}