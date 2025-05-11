import { ForwardRefExoticComponent, RefAttributes, SetStateAction, createContext, forwardRef, useCallback, useContext, useEffect, useImperativeHandle, useMemo, useRef, useState } from "react";
import { createId } from '@paralleldrive/cuid2';
import React from "react";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ChainDialogContentRef, ChainDialogRef, DialogChainPopupStack } from "./dynamic-dialog-types";
import { DynamicDialogContent } from "./dynamic-dialog-content";

export const ChainDialogContext = createContext<{
    components: DialogChainPopupStack[];

    addPopup: (name: string, contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string, beforeUpperResults?: (upperResults: any) => any) => void;
    addPopupAsFirst: (contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string) => void;

    close: (dialogId: string | null) => void;
    closeWithResult: (dialogId: string | null, data: any) => void;
    closeAll: () => void;

    setLocalState: (dialogId: string | null, data: any) => void;
}>({
    components: [],
    addPopup: () => { },
    addPopupAsFirst: () => { },

    close: (dialogId: string | null) => { },
    closeWithResult: (dialogId: string | null, result: any) => { },
    closeAll: () => { },

    setLocalState: (dialogId: string | null, data: any) => { },
});
ChainDialogContext.displayName = 'ChainDialogContext';

export const useDynamicDialogContext = () => useContext(ChainDialogContext); 

interface DynamicDialogProps {
    finish: (result: {
        data: any,
        openMode?: FormOpenMode,
        editId?: string
    }) => void;
}
export const DynamicDialog = forwardRef<ChainDialogRef, DynamicDialogProps>((props: DynamicDialogProps, ref) => {
    const [components, setComponents] = useState<Array<DialogChainPopupStack>>([]);
    
    const addPopup = useCallback((name: string, contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string, beforeUpperResults?: (upperResults: any) => any) => {
        //contentRef.current?.deactivated?.();
        setComponents((currentComponents: DialogChainPopupStack[]) => [...currentComponents, { dialogId: createId(), name, contentFC, contentFCprops, openMode, editId, beforeUpperResults }]);
    }, [setComponents]);

    const addPopupAsFirst = useCallback((contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string) => {
        setComponents([{ dialogId: createId(), name: null, contentFC, contentFCprops, openMode, editId }]);
    }, [setComponents]);

    const closeAll = useCallback(() => {
        setComponents([]);
    }, [setComponents]);

    const setLocalState = useCallback((dialogId: string | null, data: any) => {

        setComponents((currentComponents: DialogChainPopupStack[]) => currentComponents.map((i) => {
            if (i.dialogId != dialogId) {
                return i;
            }
            
            return {
                ...i,
                localState: data
            }
        }));
    }, [setComponents]);

    const close = useCallback((dialogId: string | null) => {
        const targetIndex = components.findIndex(item => item.dialogId === dialogId);
        if (targetIndex < 0) {
            return;
        }
        setComponents((currentComponents: DialogChainPopupStack[]) => currentComponents.slice(0, targetIndex));
    }, [components, setComponents]);

    const closeWithResult = useCallback((dialogId: string | null, result: any) => {
        
        if (!dialogId) {
            return;
        }

        const component = components.find(i => i.dialogId === dialogId);
        if (!component) {
            return;
        }

        const index = components.indexOf(component);
        if (index == 0) {
            setComponents([]);
            props.finish({
                data: result,
                openMode: component.openMode,
                editId: component.editId
            });
            return;
        }

        let prevComponent = { ...components[index - 1] };
        prevComponent.upperResults = {
            ...prevComponent.upperResults,
            [component.name!]: component.beforeUpperResults?.(result) ?? result
        };

        setComponents((currentComponents: DialogChainPopupStack[]) => currentComponents.slice(0, index).map((i) => {
            if (i.dialogId === prevComponent.dialogId) {
                return prevComponent;
            }
            return i;
        }));

    }, [ components, setComponents]);

    useImperativeHandle(ref, () => ({
        addPopup: (contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string) => {
            addPopupAsFirst(contentFC, contentFCprops, openMode, editId);
        }
    }));

    const provider = useMemo(() => ({
        components,
        addPopup,
        addPopupAsFirst,
        close,
        closeWithResult,
        closeAll,
        setLocalState
    }), [components, addPopup, addPopupAsFirst, close, closeWithResult, closeAll, setLocalState]);

    return (
        <ChainDialogContext.Provider value={provider}>
            <DynamicDialogContentManager />
        </ChainDialogContext.Provider>
    );
});


const DynamicDialogContentManager = () => {

    const dynaicDialogContext = useDynamicDialogContext();

    // const lastPopup = useMemo(() => {
    //     return dynaicDialogContext.components[dynaicDialogContext.components.length - 1];
    // }, [dynaicDialogContext.components]);


    // const isOpen = useMemo(() => {
    //     return !!lastPopup;
    // }, [lastPopup]);

    return (
        dynaicDialogContext.components.map((dialog, index) => {
            return (
                <DynamicDialogContent key={dialog.dialogId} dialog={dialog} dialogIndex={index} />
            );
        })
    )
};
DynamicDialog.displayName = 'DynamicDialog';

