"use client";

import { Box, FilledInput, FormControl, FormControlLabel, FormHelperText, InputLabel, Link, Stack, Switch, SwitchProps, Typography } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, get, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode, useCallback, useId, useState } from 'react';
import { UploadOneFile } from '@/shared/components/upload/upload-one-file';
import { saveFile } from '@/lib/utils/save-file';
import { useDropzone } from 'react-dropzone';
import { UploadManyFiles } from '@/shared/components/upload/upload-many-files';
import { ManyFilesPreview } from '@/shared/components/upload/components/many-files-preview';
import { FileDrop } from '@/shared/components/upload/components/file-drop';
import { BrowseButton } from '@/shared/components/upload/components/browse-button';
import AttachFileIcon from '@mui/icons-material/AttachFile'; 

interface UploadManyFilesComponentProps {
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
export function UploadManyFilesComponent(props: UploadManyFilesComponentProps) {

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
			props.onChange?.({ target: { name: props.name, value: [...value, ...acceptedFiles] } } as any); 
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

	const isFull = props.maxFiles != undefined && value?.length >= props.maxFiles;

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
				inputComponent={UploadManyFiles}
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

		<FileDrop dropzone={dropzone}>
			<ManyFilesPreview
				onRemove={onRemove} 
				firstNode={ isFull == false ? <BrowseButton dropzone={dropzone}>  
					<Box sx={{
						display: 'flex',
						borderRadius: 1,
						alignItems: 'center',
						color: 'text.disabled',
						flexDirection: 'column',
						justifyContent: 'center',
						textAlign: 'center',
					}}>
						<AttachFileIcon />
						<Typography variant="caption">Attach File</Typography>  
					</Box>
				</BrowseButton> : undefined }

				files={value}
				 />
		</FileDrop>
		
	</>)
}


interface UploadManyFilesFieldProps {
	name: string;
	placeholder?: string;
	label: string;
	disabled?: boolean;
	hide?: boolean;
	required?: boolean;

	maxFiles?: number;
}
export function UploadManyFilesField(props: UploadManyFilesFieldProps) {

	useFormDynamicContext(props.name, { disabled: props.disabled });
	const formContext = useFormContext();
	const formControllerHandler = useLiteController({ name: props.name, control: formContext.control, disabled: props.disabled });
	useFormDynamicContext(props.name, { disabled: formControllerHandler.disabled });

	if (props.hide) {
		return null;
	}

	const errors = get(formContext.formState.errors, props.name);

	return (<UploadManyFilesComponent
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={errors?.message as string}
		maxFiles={props.maxFiles}
		{...formControllerHandler}
	/>)
}