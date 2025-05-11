"use client";


import {
    Control,
    ControllerFieldState,
    ControllerRenderProps,
    FieldError,
    FieldErrorsImpl, get, Merge,
    useFormContext,
    UseFormRegister
} from "react-hook-form";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";
import { Controller } from "react-hook-form";
import dayjs from "dayjs";
import { localeConfig } from "@/config/locale-config";
import { DateTimeToolbar } from '../pickers-extensions/date-time-toolbar';
import { DateTimeTextField } from "../pickers-extensions/date-time-textfield";
import { useLiteController } from "../lite-controller";
import { useFormDynamicContext } from "@/core-features/dynamic-form2/dynamic-form";


interface DateTimeFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    min?: dayjs.Dayjs;
    max?: dayjs.Dayjs;
    errors?: string;

    value: any;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}
export function DateTimeFormComponent(props: DateTimeFormComponentProps) {
    const { name, label, placeholder, disabled, errors, min, max, value, onChange, onBlur } = props;

    return (
        <>

            <DateTimePicker
                name={name}
                label={label}
                disabled={disabled}
                views={["year", "month", "day", "hours", "minutes", "seconds"]}
                minDate={min ? dayjs(min) : undefined}
                maxDate={max ? dayjs(max) : undefined}

                value={value ?? null}
                onChange={(value) => {
                    onChange({ target: { name, value } })
                }}
                
                slots={{
                    textField: DateTimeTextField,
                    toolbar: DateTimeToolbar as any,
                }}
                slotProps={{
                    actionBar: { actions: ['clear', 'accept'] },
                    toolbar: {
                        formValue: value
                    },
                    textField: {
                        fullWidth: true,
                        error: !!errors,
                        helperText: errors,
                        inputFormat: (value: any) => value.format(localeConfig.dateTime.inputFormat),
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

        </>
    );
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
export function DateTimeFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled});
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });
    
    if(props.hide){
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<DateTimeFormComponent
            label={props.label}
            placeholder={props.placeholder}
            required={props.required}
            errors={errors?.message as string}
            min={props.min}
            max={props.max}
            {... controllerHandler}
    />)
}


declare module '@mui/x-date-pickers/internals' {
    interface ExportedBaseToolbarProps {
        formValue: any;
    }
}
