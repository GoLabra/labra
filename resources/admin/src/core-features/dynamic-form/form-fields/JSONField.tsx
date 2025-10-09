"use client";

import { TextField } from '@mui/material';
import { ChangeHandler, Control, FieldError, FieldErrorsImpl, get, Merge, RefCallBack, useFormContext, UseFormRegister } from 'react-hook-form';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';


interface JSONFormComponentProps {
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

export function JSONFormComponent(props: JSONFormComponentProps) {
    const { name, label, placeholder, disabled, errors, onChange, onBlur, inputRef, maxLength, minRows, maxRows } = props;

    return (<>

        <TextField
            fullWidth
            name={name}
            onChange={onChange}
            onBlur={onBlur}
            ref={inputRef}
            label={label}
            placeholder={placeholder}
            error={!!errors}
            slotProps={{
                htmlInput: {
                    maxLength: maxLength,
                    required: props.required
                }
            }}
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
export function JSONFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const { ref, ...registerMethods } = formContext.register(props.name, { disabled: props.disabled, required: props.required, shouldUnregister: true });
    useFormDynamicContext(props.name, { disabled: registerMethods.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<JSONFormComponent
        label={props.label}
        placeholder={props.placeholder}
        errors={errors?.message as string}
        inputRef={ref}
        
        minRows={3}
        maxRows={6}
        
        {...registerMethods}
    />)
}