import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ChainDialogContentRef, ChainDialogRef } from "./dynamic-dialog-types";
import { useRef, useCallback, ForwardRefExoticComponent, RefAttributes } from "react";

export const useDynamicDialog = () => {

    const ref = useRef<ChainDialogRef | null>(null);

    const addPopup = useCallback((contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, propsFC?: any, openMode?: FormOpenMode, editId?: string) => {
        ref.current?.addPopup(contentFC, propsFC, openMode, editId);
    }, [ref]);

    return {
        ref,
        addPopup
    }
}
