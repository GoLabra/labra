"use client";

import { Box, FilledInput, FormControl, FormControlLabel, FormHelperText, InputLabel, Link, Stack, Switch, SwitchProps, Typography } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, get, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode, useCallback, useId, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { UploadFiles } from '@/shared/components/upload/upload';


interface UploadFilesComponentProps {
	name: string;
	label: string;
	placeholder?: string;
	disabled?: boolean;
	required?: boolean;
	errors?: ReactNode;

	value: any[];
	onChange?: (event: any) => void;
	onBlur?: (event: any) => void;
	color?: SwitchProps['color'];
	maxFiles?: number;
}
export function UploadFilesComponent(props: UploadFilesComponentProps) {

	const { name, label, placeholder, disabled, errors, value, onChange, onBlur } = props;

	const dropzone = useDropzone({
		multiple: true,
		maxFiles: props.maxFiles,
		disabled: props.disabled,
		accept: { '*/*': [] },
		onDrop: (acceptedFiles, fileRejections, event) => {
			if (!acceptedFiles.length) {
				return;
			}
			props.onChange?.({ target: { name: props.name, value: [...value ?? [], ...acceptedFiles] } } as any); 
		}
	});

	const [focused, setFocused] = useState(false);
	const id = useId();

	const onRemove = useCallback((file: File | string) => {

		const oldFile = value.find(i => i === file);
		if (oldFile < 0) {
			return;
		}

		props.onChange?.({ target: { name: props.name, value: value.filter(i => i != file)  } } as any); 
	}, [value, props.onChange]);

	// const isFull = props.maxFiles != undefined && value?.length >= props.maxFiles;

	return (<>

		{/* <FormControl
			fullWidth
			variant="filled"
			focused={focused}
			onFocus={() => setFocused(true)}
			onBlur={() => setFocused(false)}
			error={!!errors}>

			<InputLabel htmlFor={id} id={`${id}-label`}>{label}</InputLabel>

			<FilledInput
				fullWidth
				inputComponent={UploadFiles}
				inputProps={{
					dropzone
				}}
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
					padding: 0
				}} />

			{errors && (<FormHelperText>{errors}</FormHelperText>)}

			
		</FormControl> */}

		<InputLabel htmlFor={id} id={`${id}-name`}>{label}</InputLabel>

		<UploadFiles dropzone={dropzone} value={value} onRemove={onRemove} maxFiles={props.maxFiles} />
		
	</>)
}


interface UploadFilesFieldProps {
	name: string;
	placeholder?: string;
	label: string;
	disabled?: boolean;
	hide?: boolean;
	required?: boolean;

	maxFiles?: number;
}
export function UploadFilesField(props: UploadFilesFieldProps) {

	useFormDynamicContext(props.name, { disabled: props.disabled });
	const formContext = useFormContext();
	const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
	useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

	if (props.hide) {
		return null;
	}

	const errors = get(formContext.formState.errors, props.name);

	return (<UploadFilesComponent
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={errors?.message as string}
		maxFiles={props.maxFiles}
		{...formControllerHandler}
	/>)
}