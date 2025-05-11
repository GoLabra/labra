import { FilledInput, FormControl, InputLabel } from "@mui/material"
import { ComponentType, ReactNode, useState } from "react";

interface ComponentFieldProps {
    componentType: ComponentType
    componentTypeProps?: any;
    id?: string;
}
export const ComponentField = (props: ComponentFieldProps) => {

    const { componentType, componentTypeProps, id } = props;

    return (<FilledInput
        inputComponent={componentType as any}
        inputProps={componentTypeProps}
        fullWidth id={id}
        sx={{
            padding: 0,
        }} />
    )
}