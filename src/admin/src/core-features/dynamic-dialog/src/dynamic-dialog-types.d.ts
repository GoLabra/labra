import { ForwardRefExoticComponent, RefAttributes } from "react";
import { FormOpenMode } from "../../dynamic-form/form-field";

export interface ChainDialogRef {
    addPopup: (contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, propsFC?: any, openMode?: FormOpenMode, editId?: string) => void;
}

export interface ChainDialogContentRef {
    deactivated?: () => void;
    enterPressed?: (e: React.KeyboardEvent<HTMLDivElement>) => void;
}

export type DialogChainPopupStack = {
    contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>;
    contentFCprops?: any;
    localState?: any;
    dialogId: string | null;
    name: string | null;
    upperResults?: Record<string, any>;
    beforeUpperResults?: (upperResults: any) => any;
    openMode?: FormOpenMode;
    editId?:string;
}