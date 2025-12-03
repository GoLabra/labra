'use client'

import { createContext, PropsWithChildren, use, useContext } from "react";
import { useSearchParams } from "next/navigation";

export const CurrentEntityName = createContext<string | undefined>(undefined);

export const useCurrentEntityNameContext = () => useContext(CurrentEntityName);

interface CurrentEntityProviderProps {
}
export const CurrentEntityProvider = (props: PropsWithChildren<CurrentEntityProviderProps>) => {

    const entityId = useSearchParams().get('e') ?? undefined;

    return (
        <CurrentEntityName.Provider value={entityId}>
            {props.children}
        </CurrentEntityName.Provider>
    )
}