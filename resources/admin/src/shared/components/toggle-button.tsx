import IconButton from '@mui/material/IconButton';
import ModeStandbyIcon from '@mui/icons-material/ModeStandby';
import FiberManualRecordOutlinedIcon from '@mui/icons-material/FiberManualRecordOutlined';
import { ReactNode, useCallback, useEffect, useState } from 'react';
import { valueFromASTUntyped } from 'graphql';
import DeleteIcon from "@mui/icons-material/Delete";
import { Button } from '@mui/material';


interface ToggleButtonProps {
    value?: boolean;
    onValueChanged?: (value: boolean) => void;
    sx?: any;
    children?: ReactNode;
    icon?: any;
}
export const ToggleButton = (props: ToggleButtonProps) => {

    const { onValueChanged, sx, children, value, icon } = props;

    const onClick = useCallback(() => {
        onValueChanged?.(!value);
    }, [onValueChanged, value]);

    return (

        <Button
            {...(value ? { variant: 'contained' } : { variant: 'outlined' })}
            onClick={onClick}
            startIcon={icon} size="small">
            {children}
        </Button>

        // <IconButton
        //     color="secondary"
        //     aria-label="add an alarm"
        //     onClick={onClick}
        //     sx={{
        //         ...(value && { color: 'primary.main' }),
        //         ...props.sx
        //     }}
        //     title="Utc now()"
        // >
        //     { children }
        // </IconButton>
    );
};
