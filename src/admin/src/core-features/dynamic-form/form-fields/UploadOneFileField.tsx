"use client";


import { Box, FilledInput, FormControl, FormControlLabel, FormHelperText, InputLabel, Link, Switch, SwitchProps } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, get, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode, useId, useState } from 'react';
import { UploadOneFile } from '@/shared/components/upload/upload-one-file';
import { saveFile } from '@/lib/utils/save-file';


interface UploadOneFileFieldProps {
	name: string;
	label: string;
	placeholder?: string;
	disabled?: boolean;
	required?: boolean;
	errors?: ReactNode;

	value: any;
	onChange?: (event: any) => void;
	onBlur?: (event: any) => void;
	color?: SwitchProps['color']
}
export function UploadOneFileComponent(props: UploadOneFileFieldProps) {

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
				inputComponent={UploadOneFile}
				placeholder={placeholder}
				disabled={disabled}
				id={id}
				name={name}
				value={value ?? ''}
				onChange={(e) => {
					onChange?.(e)
				}}
				sx={{
					overflow: 'inherit',
					padding: 0,
					width: 'fit-content',
				}} />

			{errors && (<FormHelperText>{errors}</FormHelperText>)}

			{value && (
				<>
					<Link href="#" underline="always" fontSize="small" onClick={() => {
						saveFile(props.name, value as any);
					}}>
						Dowload
					</Link>

					<Link href="#" underline="always" fontSize="small" onClick={() => {
						onChange?.({ target: { name: props.name, value: null } } as any);
					}}>
						Remove
					</Link>
				</>)}
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
export function UploadOneFileField(props: FormFieldProps) {

	useFormDynamicContext(props.name, { disabled: props.disabled });
	const formContext = useFormContext();
	const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
	useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

	if (props.hide) {
		return null;
	}

	const errors = get(formContext.formState.errors, props.name);

	return (<UploadOneFileComponent
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={errors?.message as string}
		{...formControllerHandler}
	/>)
}