import { forwardRef, memo, useEffect, useMemo, useState } from 'react';
import { DropzoneState, useDropzone } from 'react-dropzone';
import { Image } from "@/shared/components/image/image";
import { Box, InputBaseComponentProps, Stack, Typography } from '@mui/material';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import { FileDrop } from './components/file-drop';
import { OneFilePreview } from './components/one-file-preview';

export const UploadManyFiles = memo((props: InputBaseComponentProps) => {

	const { dropzone } = props;

	const hasFile = !!props.value;
	const hasError = dropzone.isDragReject;

	const renderPlaceholder = (
		<>
			{/* <Box sx={{
				height: 0,
				...(hasFile == false && { width: 200 })
			}} /> */}

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
					//position: 'absolute',
					alignItems: 'center',
					color: 'text.disabled',
					flexDirection: 'column',
					justifyContent: 'center',
					// bgcolor: 'color-mix(in srgb, var(--mui-palette-grey-500), transparent 80%)',
					transition: (theme) =>
						theme.transitions.create(['opacity'], { duration: theme.transitions.duration.shorter }),
					'&:hover': { opacity: 0.72 },
				}}
			>



				<Stack spacing={1} sx={{ textAlign: 'center' }}>
					<Box sx={{ typography: 'h6' }}>Drop or select file</Box>
					<Box sx={{ typography: 'body2', color: 'text.secondary' }}>
						Drop files here or click to
						<Box
							component="span"
							sx={{ mx: 0.5, color: 'primary.main', textDecoration: 'underline' }}
						>
							browse
						</Box>
						through your machine.
					</Box>
				</Stack>



			</Box>
		</>
	);

	return (
		<>
			<FileDrop dropzone={dropzone} sx={{width: 1}}>
				{renderPlaceholder}
			</FileDrop>
		</>
	);
});