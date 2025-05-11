"use client";

import { Control, FieldError, FieldErrorsImpl, get, Merge, useController, useFormContext, UseFormRegister } from 'react-hook-form';
import { getValueFromOptions, SelectFormComponent } from './SelectField';
import { useEffect, useMemo } from 'react';
import { Options } from '@/core-features/dynamic-form/form-field';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

interface BooleanSelectFieldProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    options: Options;
    errors?: string;

    value: any;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}

export function SingleChoiceFieldFormComponent(props: BooleanSelectFieldProps) {
    const { name, label, placeholder, disabled, errors, required, options, value, onChange, onBlur } = props;

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
    required?: boolean;
    options: Options;
    hide?: boolean;
}
export function SingleChoiceFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<SingleChoiceFieldFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        options={props.options}
        {...controllerHandler}
    />)
}
