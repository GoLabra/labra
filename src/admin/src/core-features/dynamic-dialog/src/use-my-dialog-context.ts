import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { createContext, ForwardRefExoticComponent, RefAttributes, useCallback, useContext, useMemo } from "react";
import { ChainDialogContentRef } from "./dynamic-dialog-types";
import { useDynamicDialogContext } from "./dynamic-dialog";

export const MyDialogContext = createContext<{
    dialogId: string | null;
    upperResults: any;
    localState: any;
    openMode: FormOpenMode | undefined;
    editId: string | undefined;

    addPopup: (name: string, contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string, beforeUpperResults?: (upperResults: any) => any) => void; 
}>(null!);

export const useMyDialogContext = () => {

    const dynaicDialogContext = useDynamicDialogContext();
    const myDialogContext = useContext(MyDialogContext);

    const setLocalState = useCallback((data: any) => {
        dynaicDialogContext.setLocalState(myDialogContext.dialogId, data);
    }, [dynaicDialogContext.setLocalState, myDialogContext.dialogId]);

    const close = useCallback(() => {
        dynaicDialogContext.close(myDialogContext.dialogId);
    }, [dynaicDialogContext.close, myDialogContext.dialogId]);

    const closeWithResult = useCallback((data: any) => {
        dynaicDialogContext.closeWithResult(myDialogContext.dialogId, data);
    }, [dynaicDialogContext.closeWithResult, myDialogContext.dialogId]);

    return useMemo(() => ({
        setLocalState: setLocalState,
        upperResults: myDialogContext.upperResults,
        localState: myDialogContext.localState,
        openMode: myDialogContext.openMode,
        editId: myDialogContext.editId,

        close: close,
        closeWithResults: closeWithResult,

        addPopup: myDialogContext.addPopup,
        addPopupAsFirst: dynaicDialogContext.addPopupAsFirst,
    }), [myDialogContext.upperResults, myDialogContext.localState, myDialogContext.openMode, myDialogContext.editId, setLocalState]);
}
