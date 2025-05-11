"use client";
import { Options } from '@/core-features/dynamic-form/form-field';
import { Box, Chip, FilledInput, FormControl, FormControlLabel, FormHelperText, InputLabel, MenuItem, OutlinedInput, Select, SelectChangeEvent, Stack, TextField, Theme, useTheme } from '@mui/material';
import { useId } from 'react'; 
import { get, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { selectValue } from '@/lib/utils/select-value';
import CheckIcon from '@mui/icons-material/Check';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';

export const getValuesFromOptions = (value: Array<string | number | boolean>, options: Options):
    (string | number | boolean | Record<string, any>)[] => {

    var option = getOptionsFromOptions(value, options);
    return option?.map(i => i?.value ?? null).filter(i => i != null) ?? [];
}

export const getOptionsFromOptions = (value: Array<string | number | boolean>, options: Options) => {

    if (!options) {
        return null;
    }

    if (!value || !value.length) {
        return null;
    }

    return value.map(val => options.find(option => option.value === val)).filter(Boolean);
}


const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
    PaperProps: {
        style: {
            maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
            width: 250,
        },
    },
};


interface SelectFormComponentProps {
    name: string;
    label: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    options: Options;
    errors?: string;

    value: Array<string | number | boolean>;
    onChange: (event: any) => void;
    onBlur: (event: any) => void;
}

export function TagsSelectFormComponent(props: SelectFormComponentProps) {
    const { name, label, options, placeholder, disabled, errors, value, onChange, onBlur } = props;

    const id = useId();

    return (
        <FormControl fullWidth variant="filled"
            error={!!errors}
        >
            <InputLabel htmlFor={id}>{label}</InputLabel>
            <Select
                label={label}
                name={name}
                multiple
                value={value ?? []}
                onChange={onChange}
                disabled={disabled}

                input={<FilledInput fullWidth id={id} sx={{ padding: '4px' }} />}
                renderValue={(selected: Array<string | number | boolean>) => {

                    const selectedOptions = getOptionsFromOptions(selected, options);

                    return <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selectedOptions?.map((selectedOption) => (
                            <Chip key={selectedOption!.label} label={selectedOption!.label} />
                        ))}
                    </Box>

                }}
                MenuProps={MenuProps}
            >
                {options.map((option) => (
                    <MenuItem
                        key={option.label}
                        value={option.value as any}
                    >
                        {option.label}
                    </MenuItem>
                ))}
            </Select>

            {errors && (<FormHelperText>{errors}</FormHelperText>)}
        </FormControl>
    )
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
export function TagsSelectFormField(props: FormFieldProps) {
    const formContext = useFormContext();
    const controllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
    useFormDynamicContext(props.name, { disabled: controllerHandler.disabled });

    if (props.hide) {
        return null;
    }

    const errors = get(formContext.formState.errors, props.name);

    return (<TagsSelectFormComponent
        label={props.label}
        placeholder={props.placeholder}
        required={props.required}
        errors={errors?.message as string}
        options={props.options}
        {...controllerHandler}
    />)
}
