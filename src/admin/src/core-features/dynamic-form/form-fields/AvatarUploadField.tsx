"use client";


import { Box, FormControlLabel, FormHelperText, Switch } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode } from 'react';
import { UploadAvatar } from '@/shared/components/upload/upload-avatar';

interface FormFieldProps {
    name: string;
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function AvatarUploadField(props: FormFieldProps) {

    useFormDynamicContext(props.name, { disabled: props.disabled });
    const formContext = useFormContext();
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = formContext.formState.errors[props.name]?.message as string;

    return (<>

        <UploadAvatar
            value={formControllerHandler.value}
            onDrop={(acceptedFiles) => {
                formControllerHandler.onChange({ target: { name: props.name, value: acceptedFiles[0] } })
            }}
            disabled={formControllerHandler.disabled}
            errors={formContext.formState.errors[props.name]?.message as string}
        />

        {!!errors && (
            <FormHelperText error sx={{ px: 2, textAlign: 'center' }}>
                {errors}
            </FormHelperText>
        )}
    </>)
}