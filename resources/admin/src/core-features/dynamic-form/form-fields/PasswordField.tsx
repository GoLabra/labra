"use client";

import { TextField } from '@mui/material';
import { ChangeHandler, get, RefCallBack, useFormContext } from 'react-hook-form';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

interface PasswordFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean
    errors?: string;

    onChange?: ChangeHandler;
    onBlur?: ChangeHandler;
    inputRef?: RefCallBack;
    maxLength?: number;
}

export function PasswordFormComponent(props: PasswordFormComponentProps) {
    const { name, label, placeholder, disabled, errors, onChange, onBlur, inputRef, maxLength } = props;

    return (<>

        <TextField
            fullWidth
            name={name}
            onChange={onChange}
            onBlur={onBlur}
            ref={inputRef}
            type="password"
            label={label}
            placeholder={placeholder}
            error={!!errors}
            slotProps={{
                htmlInput: {
                    maxLength: maxLength,
                    required: props.required
                }
            }}
            disabled={disabled}
            helperText={errors ?? null}
        />

    </>)
}
PasswordFormComponent.displayName = 'RichTextInput';



interface PasswordFormField {
    name: string;
    placeholder?: string;
    label: string;
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function PasswordFormField(props: PasswordFormField) {

    const formContext = useFormContext();
    const { ref, ...registerMethods } = formContext.register(props.name, { disabled: props.disabled, required: props.required, shouldUnregister: true, maxLength: 256 });
    useFormDynamicContext(props.name, { disabled: registerMethods.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<PasswordFormComponent
        label={props.label}
        placeholder={props.placeholder}
        errors={errors?.message as string}
        inputRef={ref}
        {...registerMethods}
    />)
}