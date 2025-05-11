"use client";

import { TextField } from '@mui/material';
import {  ChangeHandler, get, RefCallBack, useFormContext } from 'react-hook-form';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';


interface NumberFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    min?: string | number;
    max?: string | number;
    required?: boolean;
    errors?: string;

    onChange?: ChangeHandler;
    onBlur?: ChangeHandler;
    inputRef?: RefCallBack;
}

export function NumberFormComponent(props: NumberFormComponentProps) {
    const { name, label, placeholder, disabled, errors, min, max, onChange, onBlur, inputRef } = props;
    
    var minValue = min ?? -2147483648;
    var maxValue = max ?? 2147483647;    

    return (<>
            <TextField
                fullWidth
                type='number'
                name={name}
                onChange={onChange}
                onBlur={onBlur}
                ref={inputRef}
                label={label}
                placeholder={placeholder}
                error={!!errors}
                slotProps={{
                    htmlInput: {
                        min: minValue,
                        max: maxValue,
                    }
                }}
                disabled={disabled}
                helperText={errors}
            />
        
    </>)
}

interface FormFieldProps {
    name: string;
    placeholder?: string;
    label: string;
    min?: string | number;
    max?: string | number
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function NumberFormField(props: FormFieldProps) {
    
    const formContext = useFormContext();
    const {ref, ...registerMethods} = formContext.register(props.name, { disabled: props.disabled, required: props.required, shouldUnregister: true, min: props.min, max: props.max }); 
    useFormDynamicContext(props.name, { disabled: registerMethods.disabled });
    
    const errors = get(formContext.formState.errors, props.name);
    
    return (<NumberFormComponent
            label={props.label}
            placeholder={props.placeholder}
            errors={errors?.message as string}
            min={props.min}
            max={props.max}
            inputRef={ref}
            {... registerMethods}
    />)
}