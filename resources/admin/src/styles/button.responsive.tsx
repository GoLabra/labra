import styled from "@emotion/styled";
import React, { forwardRef, type ComponentType } from 'react';
import { Button, ButtonProps, useMediaQuery, useTheme } from "@mui/material";


const StyledResponsiveButton = styled(Button, {
    shouldForwardProp: (prop) => prop !== 'mobileProps',
})<ResponsiveButtonProps>(({ theme }) => ({
    fontSize: (theme as any).typography.pxToRem(14),
    minWidth: 'auto',
    [(theme as any).breakpoints.down('md')]: {
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
}));

export type ResponsiveButtonProps = ButtonProps & {
    mobileProps?: Partial<ButtonProps>;
};

export const ResponsiveButton = forwardRef<HTMLButtonElement, ResponsiveButtonProps>((props, ref) => {
    const theme = useTheme();
    const isMobile = useMediaQuery((theme as any).breakpoints.down('md'));

    const { mobileProps, ...rest } = props;

    if (isMobile) {
        const mergedMobileProps: ButtonProps = {
            ...rest,
            ...(mobileProps ?? {}),
        } as ButtonProps;
        return <StyledResponsiveButton ref={ref} {...mergedMobileProps} />;
    }

    return <StyledResponsiveButton ref={ref} {...rest} />;
});
ResponsiveButton.displayName = 'ResponsiveButton';

// (ResponsiveButton as unknown as ComponentType<ButtonProps>).defaultProps = {
//     variant: 'text',
//     size: 'medium',
//     color: 'secondary',
// };


