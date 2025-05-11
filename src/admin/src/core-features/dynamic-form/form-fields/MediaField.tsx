"use client";

import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { TextField } from '@mui/material';
import { FieldError, FieldErrorsImpl, get, Merge, useFormContext, UseFormRegister } from 'react-hook-form';
import { useLiteController } from '../lite-controller';


interface MediaFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: string;
    value: any;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}

export function MediaFormComponent(props: MediaFormComponentProps) {
    const { name, label, placeholder, disabled, errors, value, onChange, onBlur } = props;

    return (<>

        MEDIA

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
export function MediaFormField(props: FormFieldProps) {
    useFormDynamicContext(props.name, { disabled: props.disabled });
    const formContext = useFormContext();
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<MediaFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        {... formControllerHandler}
    />)
}