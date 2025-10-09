"use client";

import { TextField } from '@mui/material';
import { ChangeHandler, FieldError, FieldErrorsImpl, get, Merge, RefCallBack, useFormContext, UseFormRegister } from 'react-hook-form';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';


interface TextLongFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: string;

    onChange?: ChangeHandler;
    onBlur?: ChangeHandler;
    inputRef?: RefCallBack;
    maxLength?: number;
    minRows?: number;
    maxRows?: number;
}

export function TextLongFormComponent(props: TextLongFormComponentProps) {
    const { name, label, placeholder, disabled, errors, onChange, onBlur, inputRef, maxLength, minRows, maxRows } = props;

    return (<>

        <TextField
            fullWidth
            name={name}
            onChange={onChange}
            onBlur={onBlur}
            ref={inputRef}
            label={label}
            disabled={disabled}
            placeholder={placeholder}
            error={!!errors}
            inputProps={{ maxLength: maxLength }}
            helperText={errors}
            multiline
            minRows={minRows}
            maxRows={maxRows}
        />

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
export function TextLongFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const { ref, ...registerMethods } = formContext.register(props.name, { disabled: props.disabled, required: props.required, shouldUnregister: true, maxLength: 256 });
    useFormDynamicContext(props.name, { disabled: registerMethods.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<TextLongFormComponent
        label={props.label}
        placeholder={props.placeholder}
        errors={errors?.message as string}

        minRows={3}
        maxRows={6}

        inputRef={ref}
        {...registerMethods}

    />)
}
