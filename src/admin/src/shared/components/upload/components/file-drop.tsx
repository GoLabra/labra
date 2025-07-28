import { Box, SxProps, Theme } from "@mui/material";
import { PropsWithChildren } from "react";
import { Accept, DropzoneOptions, DropzoneState, useDropzone } from 'react-dropzone';

interface FileDropProps {
	name: string;
	dropzone: DropzoneState
	disabled?: boolean,
	sx?: SxProps<Theme>
}
export const FileDrop = (props: PropsWithChildren<FileDropProps>) => {

	const { dropzone } = props;
	const hasError = dropzone.isDragReject;

	const { onClick, ...other} = dropzone.getRootProps();

	return (
		<Box
			{...other}
			sx={{
				outline: 'none !important',
				// width: 'fit-content',
				// cursor: 'pointer',
				...(dropzone.isDragActive && { opacity: 0.60 }),
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
			<input name={props.name} {...dropzone.getInputProps()} />

			{props.children}
		</Box>
	);
}