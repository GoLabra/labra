"use client";

import { Box, FilledInput, FormControl, FormControlLabel, FormHelperText, InputLabel, Link, Stack, Switch, SwitchProps, ToggleButton, ToggleButtonGroup, Typography } from '@mui/material';
import { Control, FieldError, FieldErrorsImpl, Merge, UseFormRegister, get, useController, useFormContext } from 'react-hook-form';
import { useLiteController } from '../lite-controller';
import { useFormDynamicContext } from '@/core-features/dynamic-form2/dynamic-form';
import { ReactNode, useCallback, useId, useState } from 'react';
import { DropEvent, DropzoneOptions, FileRejection, useDropzone } from 'react-dropzone';
import { UploadFiles } from '@/shared/components/upload/upload';
import ViewListIcon from '@mui/icons-material/ViewList';
import ViewModuleIcon from '@mui/icons-material/ViewModule';

interface UploadFilesBaseFieldProps {
	name: string;
	label: string;
	placeholder?: string;
	disabled?: boolean;
	required?: boolean;
	errors?: ReactNode;

	maxFiles?: number;

	value: (File | string)[];
	onChange?: (event: any) => void;
	onBlur?: (event: any) => void;

	onDrop?: DropzoneOptions['onDrop'];
	onRemove?: (file: File | string) => void;
}
export function   UploadFilesBaseField(props: UploadFilesBaseFieldProps) {

	const { name, label, placeholder, disabled, errors, value, onChange, onBlur } = props;

	const dropzone = useDropzone({
		multiple: true,
		maxFiles: props.maxFiles,
		disabled: props.disabled,
		accept: { '*/*': [] },
		onDrop: (acceptedFiles, fileRejections, event) => {
			
			props.onDrop?.(acceptedFiles, fileRejections, event);
			
			if (!acceptedFiles.length) {
				return;
			}
			
			if(!props.onChange){
				return;
			}

			props.onChange?.({ target: { name: props.name, value: [...value ?? [], ...acceptedFiles] } } as any);
		}
	});

	const [focused, setFocused] = useState(false);
	const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

	const id = useId();

	const onRemove = useCallback((file: File | string) => {

		props.onRemove?.(file);
		
		if(!props.onChange){
			return;
		}

		const oldFile = value.find(i => i === file);
		if (!oldFile) {
			return;
		}

		props.onChange?.({ target: { name: props.name, value: value.filter(i => i != file) } } as any);
	}, [value, props.onChange]);

	const handleChange = (event: React.MouseEvent<HTMLElement>, nextView: 'grid' | 'list') => {
		setViewMode(nextView);
	};

	return (<>

		<FormControl
			fullWidth
			variant="filled"
			focused={focused}
			onFocus={() => setFocused(true)}
			onBlur={() => setFocused(false)}
			error={!!errors}>

			<InputLabel htmlFor={id} id={`${id}-name`} shrink={true}>
				<Stack direction="row" alignItems="center" justifyContent="space-between" gap={1}>
					{label}

					<ToggleButtonGroup
						value={viewMode}
						exclusive
						onChange={handleChange}
					>
						<ToggleButton value="list" aria-label="list" size="small" sx={{ padding: 0.5 }}>
							<ViewListIcon />
						</ToggleButton>
						<ToggleButton value="grid" aria-label="grid" size="small" sx={{ padding: 0.5 }}>
							<ViewModuleIcon />
						</ToggleButton>
					</ToggleButtonGroup>
				</Stack>
			</InputLabel>

			<UploadFiles name={name} dropzone={dropzone} value={value} onRemove={onRemove} maxFiles={props.maxFiles} viewMode={viewMode} />

			{errors && (<FormHelperText>{errors}</FormHelperText>)}

		</FormControl>

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

	return (<UploadFilesBaseField
		label={props.label}
		placeholder={props.placeholder}
		required={props.required}
		errors={errors?.message as string}
		maxFiles={props.maxFiles}
		{...formControllerHandler}
	/>)
}