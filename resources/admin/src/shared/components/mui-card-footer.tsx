import { Box, styled } from "@mui/material"
import { ReactNode } from "react";

export const MuiCardFooterRoot = styled(Box, {
    name: 'MuiCardFooter', // The component name
    slot: 'root', // The slot name
})(({ theme }) => ({
    
}));

interface MuiCardFooterProps {
    children: ReactNode;
    root?: any
}
export const MuiCardFooter = (props: MuiCardFooterProps) => {

    const { children, root } = props

    return (<MuiCardFooterRoot className="MuiCardFooter" {...root}>
        {children}
    </MuiCardFooterRoot>)
}