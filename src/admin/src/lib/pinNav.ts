"use client";

import { makeVar, useReactiveVar } from "@apollo/client";
import { useCallback, useMemo } from "react";

export const pinNavVar = makeVar<boolean>(true);

export const usePinNav = () => {
    const pinNav = useReactiveVar(pinNavVar);

    const togglePinNav = useCallback(() => {
        pinNavVar(!pinNav)
    }, [pinNav]);

    return useMemo(() => ({
        pinNav,
        togglePinNav
    }), [pinNav, togglePinNav]);
}
