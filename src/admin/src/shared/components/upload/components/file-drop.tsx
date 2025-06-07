import { Box, SxProps, Theme } from "@mui/material";
import { PropsWithChildren } from "react";
import { Accept, DropzoneOptions, DropzoneState, useDropzone } from 'react-dropzone';

interface FileDropProps {
	dropzone: DropzoneState
	disabled?: boolean,
	sx?: SxProps<Theme>
}
export const FileDrop = (props: PropsWithChildren<FileDropProps>) => {

	const { dropzone } = props;
	const hasError = dropzone.isDragReject;

	return (
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
				// ...(hasFile && {
				// 	...(hasError && {
				// 		bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
				// 	}),
				// 	'&:hover .upload-placeholder': { opacity: 1 },
				// }),
				...props.sx,
			}}
		>
			<input {...dropzone.getInputProps()} />

			{props.children}
		</Box>
	);
}