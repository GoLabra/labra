import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { Dialog, Slide, Theme, useMediaQuery } from "@mui/material";
import { forwardRef, ForwardRefExoticComponent, RefAttributes, useCallback, useMemo, useRef } from "react";
import { ChainDialogContentRef, DialogChainPopupStack } from "./dynamic-dialog-types";
import { useDynamicDialogContext } from "./dynamic-dialog";
import { TransitionProps } from "@mui/material/transitions";
import { MyDialogContext } from "./use-my-dialog-context";
import { Key, useKeyHandler } from "@/shared/components/key-handler";

interface DynamicDialogContentProps {
    dialog: DialogChainPopupStack;
    dialogIndex: number;
}
export const DynamicDialogContent = (props: DynamicDialogContentProps) => {

    const dynaicDialogContext = useDynamicDialogContext();
    const ContentFC = useMemo(() => props.dialog.contentFC, [props.dialog]);
    const contentRef = useRef<ChainDialogContentRef>();
    const keyHandler = useKeyHandler({ keyToHandle: Key.Enter, modifiers:["ctrl"], onPressed: (e) => contentRef.current?.enterPressed?.(e) });
    const mdDown = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

    const dialogStyle = useMemo(() => {
        if(mdDown){
            return { };
        }
        return {
            marginLeft: props.dialogIndex * 3
        }
    }, [props.dialogIndex, mdDown]);

    const addPopup = useCallback((name: string, contentFC: ForwardRefExoticComponent<RefAttributes<ChainDialogContentRef> & any>, contentFCprops?: any, openMode?: FormOpenMode, editId?: string, beforeUpperResults?: (upperResults: any) => any) => {
        contentRef.current?.deactivated?.();
        dynaicDialogContext.addPopup(name, contentFC, contentFCprops, openMode, editId, beforeUpperResults);
    }, [dynaicDialogContext.addPopup]);

    const provider = useMemo(() => ({
        dialogId: props.dialog.dialogId,
        upperResults: props.dialog.upperResults,
        localState: props.dialog.localState,
        openMode: props.dialog.openMode,
        editId: props.dialog.editId,
        addPopup
    }), [props.dialog.dialogId, props.dialog.upperResults, props.dialog.localState, props.dialog.openMode, props.dialog.editId, addPopup]);

    return (
        <MyDialogContext.Provider value={provider}>

            <Dialog
                ref={keyHandler.setRef}
                keepMounted={false}
                onClose={() => dynaicDialogContext.close(props.dialog.dialogId)}
                TransitionComponent={Transition}
                open={true}
                fullWidth
                maxWidth="sm"

                PaperProps={{
                    sx: {
                        overflow: 'visible',
                        ...dialogStyle
                    }
                }}
            >
                <ContentFC
                    ref={contentRef}
                    key={`contentFC${props.dialog.dialogId}`}
                    {...props.dialog.contentFCprops} />

            </Dialog>

        </MyDialogContext.Provider>
    )
}


const Transition = forwardRef(function Transition(
    props: TransitionProps & {
        children: React.ReactElement<any, any>;
    },
    ref: React.Ref<unknown>,
) {
    return <Slide direction="down" ref={ref} {...props} />;
});

