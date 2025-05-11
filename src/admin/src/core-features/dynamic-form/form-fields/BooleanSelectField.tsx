"use client";

import { Control, FieldError, FieldErrorsImpl, get, Merge, useFormContext, UseFormRegister } from 'react-hook-form';
import { SelectFormComponent } from './SelectField';
import { useMemo } from 'react';
import { Options } from '@/core-features/dynamic-form/form-field';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

interface BooleanSelectFieldProps {
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

export function BooleanSelectFieldFormComponent(props: BooleanSelectFieldProps) {
    const { name, label, placeholder, disabled, errors, required, value, onChange, onBlur } = props;

    const options = useMemo((): Options => {
        return [{
                label: 'True',
                value: true
            },
            {
                label: 'False',
                value: false
            }]
    }, []);

    return (
        <>
            <SelectFormComponent
                name={name}
                label={label}
                placeholder={placeholder}
                disabled={disabled}
                errors={errors}
                required={required}
                options={options}
                value={value}
                onChange={onChange}
                onBlur={onBlur}>

            </SelectFormComponent>
        </>
    );
}


interface FormFieldProps {
    name: string;
    placeholder?: string;
    label: string;
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function BooleanSelectFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<BooleanSelectFieldFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        {...controllerHandler}
    />)
}