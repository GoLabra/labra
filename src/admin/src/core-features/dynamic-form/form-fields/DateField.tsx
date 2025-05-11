"use client";

import { TextField } from '@mui/material';
import { Control, Controller, FieldError, FieldErrorsImpl, get, Merge, useFormContext, UseFormRegister } from 'react-hook-form';
import { DatePicker, DatePickerToolbar } from '@mui/x-date-pickers/DatePicker';
import dayjs from 'dayjs';
import { localeConfig } from '@/config/locale-config';
import { DateTimeToolbar } from '../pickers-extensions/date-time-toolbar';
import { DateTimeTextField } from '../pickers-extensions/date-time-textfield';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

interface DateFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean
    min?: dayjs.Dayjs;
    max?: dayjs.Dayjs;
    errors?: string;

    value: any;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}
export function DateFormComponent(props: DateFormComponentProps) {
    const { name, label, placeholder, disabled, required, errors, min, max, value, onChange, onBlur } = props;

    return (<>

        <DatePicker
            name={name}
            label={label}
            disabled={disabled}
            minDate={min ? dayjs(min) : undefined}
            maxDate={max ? dayjs(max) : undefined}

            value={value ?? null}
            onChange={(value) => { onChange({ target: { name, value } }) }}

            slots={{
                textField: DateTimeTextField,
                toolbar: DateTimeToolbar as any
            }}

            slotProps={{
                actionBar: { actions: ['clear'] },
                toolbar: {
                    formValue: value,
                },
                textField: {
                    fullWidth: true,
                    error: !!errors,
                    helperText: errors,
                    inputFormat: (value: any) => value.format(localeConfig.date.inputFormat),
                    formValue: value
                } as any,
                popper: {
                    placement: "bottom-end",
                    popperOptions: {
                        modifiers: [
                            {
                                name: 'preventOverflow',
                                enabled: true,
                                options: {
                                    altAxis: true,
                                    altBoundary: false,
                                    tether: false,
                                    rootBoundary: 'viewport',
                                    padding: 0,
                                },
                            },
                        ]
                    }
                }
            }}
        />

    </>)
}

interface FormFieldProps {
    name: string;
    placeholder?: string;
    label: string;
    min?: dayjs.Dayjs;
    max?: dayjs.Dayjs;
    disabled?: boolean;
    hide?: boolean;
    required?: boolean;
}
export function DateFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<DateFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        min={props.min}
        max={props.max}
        {...controllerHandler} />
    )
}



declare module '@mui/x-date-pickers/internals' {
    interface ExportedBaseToolbarProps {
        formValue: any;
    }
}
