"use client";


import { Box, FormControlLabel, Switch } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, get, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode } from 'react';

interface BooleanFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: ReactNode;

    value: any;
    onChange?: (event: any) => void;
    onBlur?: (event: any) => void;
}
export function BooleanFormComponent(props: BooleanFormComponentProps) {
    const { name, label, placeholder, disabled, errors, value, onChange, onBlur } = props;


    return (<>
        <Box
            border={0} padding={0}>
            <FormControlLabel
                name={name}
                control={<Switch checked={value ?? false} onChange={(_, checked: boolean) => {
                    onChange?.({ target: { name, value: checked } })
                }
                } />
                }
                labelPlacement="start"
                label={label}
                disabled={disabled}

                sx={{
                    margin: 0,
                    padding: 0
                }} />
        </Box>
    </>)
}

interface FormFieldProps {
    name: string;
    placeholder?: string;
    label: string;
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function BooleanFormField(props: FormFieldProps) {

    useFormDynamicContext(props.name, { disabled: props.disabled });
    const formContext = useFormContext();
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<BooleanFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        {...formControllerHandler}
    />)
}