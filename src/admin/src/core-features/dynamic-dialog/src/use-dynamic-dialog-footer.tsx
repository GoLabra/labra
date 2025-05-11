import { ReactNode, useMemo } from "react";
import { useMyDialogContext } from "./use-my-dialog-context";
import { DialogActions, Divider } from "@mui/material";

export const DynamicDialogFooter = ({ children }: { children?: ReactNode }) => {

    return (
        <>
            <Divider />
            <DialogActions>
                {children}
            </DialogActions>
        </>
    )
}
