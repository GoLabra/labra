import styled from "@emotion/styled";
import { Button, ButtonProps } from "@mui/material";

export const ResponsiveButton = styled(Button)<ButtonProps>(({ theme }) => ({

    fontSize: (<any>theme).typography.pxToRem(14),
    minWidth: 'auto',

    [(<any>theme).breakpoints.down('md')]: {
        minWidth: 32,
        paddingLeft: 8,
        paddingRight: 8,
        '& .MuiButton-startIcon': {
            margin: 0,
        },
        '& .button-text': {
            display: 'none',
        },
    },
})) as typeof Button;