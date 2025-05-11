"use client";

import { Button, FilledInput, FormControl, FormHelperText, InputLabel } from '@mui/material';
import { get, useFormContext, } from 'react-hook-form';
import { forwardRef, ReactNode, useCallback, useImperativeHandle, useRef, useState } from 'react';
import { useLiteController } from '../lite-controller';
import useId from '@mui/material/utils/useId';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { TiptapEditor } from '@/shared/components/tiptap/tiptapEditor';

interface RichTextFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    errors?: ReactNode;

    value: any;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}

export function RichTextFormComponent(props: RichTextFormComponentProps) {
    const { name, label, placeholder, disabled, errors, value, onChange, onBlur } = props;

    const [focused, setFocused] = useState(false);
    const id = useId();

    return (<>

        <FormControl
            fullWidth
            variant="filled"
            focused={focused}
            onFocus={() => setFocused(true)}
            onBlur={() => setFocused(false)}
            error={!!errors}>

            <InputLabel htmlFor={id} id={`${id}-label`}>{label}</InputLabel>

            <FilledInput
                fullWidth 
                inputComponent={TiptapEditor}
                placeholder={placeholder}
                disabled={disabled}
                id={id}
                name={name}
                value={value ?? ''}
                onChange={(e) => {
                    onChange(e)
                }}
                sx={{
                    overflow: 'inherit',
                    padding: 0,
                }} />

            {errors && (<FormHelperText>{errors}</FormHelperText>)}
        </FormControl>
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
export function RichTextFormField(props: FormFieldProps) {

    useFormDynamicContext(props.name, { disabled: props.disabled });
    const formContext = useFormContext();
    const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<RichTextFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        {...formControllerHandler}
    />)
}