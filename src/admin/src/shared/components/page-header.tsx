import { styled, alpha} from "@mui/material";
import { Box } from "@mui/material";

export const PageHeader = styled(Box)(
    ({ theme }) => ({
        position: "sticky",
        top: 64,
        zIndex: 10,
        backgroundColor: 'var(--mui-palette-background-default)',
        padding: '8px 0',
        
        // this is to hide the shadow of the body (left and right)
        boxShadow: `-3px 0 0px 0px  var(--mui-palette-background-default), 6px 0 0px 0px var(--mui-palette-background-default)`,
    })
);