import { Button, ButtonProps, styled } from "@mui/material";

interface IconTextButtonProps extends ButtonProps {
    iconOpacity?: boolean;
}

export const IconTextButton = styled(Button, {
    shouldForwardProp: (prop) => prop !== 'iconOpacity'
})<IconTextButtonProps>(
    ({ theme, iconOpacity = false }) => ({
        '& .MuiButton-startIcon, & .MuiButton-endIcon': {
            opacity: iconOpacity ? 0.5 : 1
        }
    })
);