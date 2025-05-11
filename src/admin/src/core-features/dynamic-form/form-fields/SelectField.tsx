"use client";
import { Options } from '@/core-features/dynamic-form/form-field';
import { IconButton, InputAdornment, MenuItem, TextField as MuiTextField, styled } from '@mui/material';
import { get, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { selectValue } from '@/lib/utils/select-value';
import { ClearIcon } from '@mui/x-date-pickers/icons';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

const TextField = styled(MuiTextField)(({ theme }) => ({
    '.MuiInputAdornment-root': {
        position: 'absolute',
        right: '32px',
        display: 'none',
    },
    '&:hover .MuiInputAdornment-root': {
        display: 'flex',
    }
}));

export const getValueFromOptions = (value: string | number | boolean | Record<string, any>, options: Options) => {
    var option = getOptionFromOptions(value, options);
    return option?.value;
}

export const getOptionFromOptions = (value: string | number | boolean | Record<string, any>, options: Options) => {
    if (!options) {
        return null;
    }
    return options?.find(i => i.value == value);
}

interface SelectFormComponentProps {
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
export function SelectFormComponent(props: SelectFormComponentProps) {
    const { name, label, options, placeholder, disabled, errors, value, onChange, onBlur } = props;

    const optionValue = getValueFromOptions(value, options) ?? null;

    return (
        <TextField
            fullWidth
            label={label}
            name={name}
            disabled={disabled}
            value={optionValue ?? ''}
            onChange={onChange}
            onBlur={onBlur}
            error={!!errors}
            helperText={errors}
            select

            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <IconButton size="small" onClick={() => onChange({ target: { name, value: null } })}>
                            <ClearIcon fontSize="inherit" />
                        </IconButton>
                    </InputAdornment>)

            }}>
            {options.map((option) => (
                <MenuItem
                    key={option.label}
                    value={option.value as any}
                >
                    {option.label}
                </MenuItem>
            ))}
        </TextField>)
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
export function SelectFormField(props: FormFieldProps) {

    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });

    if (props.hide) {
        return null;
    }
    
    const errors = get(formContext.formState.errors, props.name);
    
    return (<SelectFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        options={props.options}
        {...controllerHandler}
    />)
}