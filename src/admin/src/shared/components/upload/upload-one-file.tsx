import { forwardRef, memo, useEffect, useMemo, useState } from 'react';
import { DropzoneState, useDropzone } from 'react-dropzone';
import { Image } from "@/shared/components/image/image";
import { Box, InputBaseComponentProps, Stack, Typography } from '@mui/material';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import { FileDrop } from './components/file-drop';
import { OneFilePreview } from './components/one-file-preview';


export const UploadOneFile = memo((props: InputBaseComponentProps) => {

	const { dropzone } = props;
		
	const hasFile = !!props.value;
	const hasError = dropzone.isDragReject;

	const renderPlaceholder = (
		<>
			<Box sx={{
				height: 0,
				...(hasFile == false && { width: 200 })
			}} />

			<Box
				className="upload-placeholder"
				sx={{
					top: 0,
					gap: 1,
					left: 0,
					width: 1,
					height: 1,
					zIndex: 9,
					display: 'flex',
					borderRadius: 1,
					position: 'absolute',
					alignItems: 'center',
					color: 'text.disabled',
					flexDirection: 'column',
					justifyContent: 'center',
					// bgcolor: 'color-mix(in srgb, var(--mui-palette-grey-500), transparent 80%)',
					transition: (theme) =>
						theme.transitions.create(['opacity'], { duration: theme.transitions.duration.shorter }),
					'&:hover': { opacity: 0.72 },
					...(hasError && {
						color: 'error.main',
						bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
					}),
					...(hasFile && {
						zIndex: 9,
						opacity: 0,
						color: 'common.white',
						bgcolor: 'color-mix(in srgb, var(--mui-palette-grey-900), transparent 50%)',
					}),
				}}
			>
				<CameraAltIcon />
				<Typography variant="caption">{hasFile ? 'Change file' : 'Upload file'}</Typography>
			</Box>
		</>
	);

	return (
		<>
			<FileDrop dropzone={dropzone}>
				{hasFile && <OneFilePreview file={props.value} />}
				{renderPlaceholder}
			</FileDrop>
		</>
	);
});