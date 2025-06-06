import { forwardRef, useEffect, useMemo, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { Image } from "@/shared/components/image/image";
import { Box, InputBaseComponentProps, Stack, Typography } from '@mui/material';
import { RejectionFiles } from './rejection-files';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import { isImageFile } from '../is-image-file';
import { FileThumbnail } from '../file-thumbnail/file-thumbnail';

interface FilePreviewProps {
	file?: File | string;
}
const FilePreview = (props: FilePreviewProps) => {

	const isImage = useMemo(() => isImageFile(props.file), [props.file]);

	if (!props.file) {
		return null;
	}

	if (isImage) {

		const imgSource = typeof props.file === 'string' ? props.file : URL.createObjectURL(props.file);

		return (<Image alt="avatar" src={imgSource} sx={{ width: 1, height: 1 }} />)
	}

	return (
		<Box
			sx={{
				top: 0,
				gap: 1,
				left: 0,
				width: 200,
				height: 1,
				display: 'flex',
				borderRadius: 1,
				position: 'relative',
				alignItems: 'center',
				flexDirection: 'column',
				justifyContent: 'center',
			}}
		>
			<Stack direction="row" spacing={1} alignItems="center">

				<FileThumbnail file={props.file} />
				<Stack direction="column" spacing={1}>
					{/* {props.value} */}
				</Stack>
			</Stack>
		</Box>
	)
}

export const UploadOneFile = forwardRef((props: InputBaseComponentProps, ref) => {

	const dropzone = useDropzone({
		multiple: false,
		disabled: props.disabled,
		accept: { '*/*': [] },
		onDrop: (acceptedFiles, fileRejections, event) => {
			if (!acceptedFiles.length) {
				return;
			}
			props.onChange?.({ target: { name: props.name, value: acceptedFiles[0] } } as any);
		}
	});

	const hasFile = !!props.value;
	const hasError = dropzone.isDragReject || !!props.errors;

	// const [preview, setPreview] = useState('');

	// const renderPreview = hasFile && (
	// 	<>
	// 		{isImageFile(props.value) && (
	// 			<Image alt="avatar" src={preview} sx={{ width: 1, height: 1, borderRadius: 1 }} />
	// 		)}
	// 		{isImageFile(props.value) == false && (
	// 			<Box
	// 				sx={{
	// 					top: 0,
	// 					gap: 1,
	// 					left: 0,
	// 					width: 200,
	// 					height: 1,
	// 					display: 'flex',
	// 					borderRadius: 1,
	// 					position: 'relative',
	// 					alignItems: 'center',
	// 					flexDirection: 'column',
	// 					justifyContent: 'center',
	// 				}}
	// 			>
	// 				<Stack direction="row" spacing={1} alignItems="center">

	// 					<FileThumbnail file={props.value} />
	// 					<Stack direction="column" spacing={1}>
	// 						{props.value}
	// 					</Stack>
	// 				</Stack>
	// 			</Box>
	// 		)}
	// 	</>
	// );

	// useEffect(() => {
	// 	if (typeof props.value === 'string') {
	// 		setPreview(props.value);
	// 	} else if (props.value instanceof File) {
	// 		setPreview(URL.createObjectURL(props.value));
	// 	}
	// }, [props.value]);

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

	const renderContent = (
		<Box
			sx={{
				width: 'fit-content',
				height: 1,
				overflow: 'hidden',
				borderRadius: 1,
				position: 'relative',
			}}
		>
			{hasFile && <FilePreview file={props.value} />}
			{renderPlaceholder}
		</Box>
	);

	return (
		<>
			<Box
				{...dropzone.getRootProps()}
				sx={{
					outline: 'none !important',
					padding: 0.5,
					width: 'fit-content',
					height: 124,

					cursor: 'pointer',
					...(dropzone.isDragActive && { opacity: 0.72 }),
					...(props.disabled && { opacity: 0.48, pointerEvents: 'none' }),
					...(hasError && { borderColor: 'error.main' }),
					...(hasFile && {
						...(hasError && {
							bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
						}),
						'&:hover .upload-placeholder': { opacity: 1 },
					}),
					//...props.sx,
				}}
			>
				<input {...dropzone.getInputProps()} />

				{renderContent}
			</Box>

			{props.helperText && props.helperText}

			<RejectionFiles files={dropzone.fileRejections} />
		</>
	);
});
