import { Box, Button } from "@mui/material";
import { PropsWithChildren } from "react";
import { DropzoneState } from "react-dropzone";

interface BrowseButtonProps {
	dropzone: DropzoneState
}
export const BrowseButton = (props: PropsWithChildren<BrowseButtonProps>) => {

	const { dropzone } = props;
	const { onClick } = dropzone.getRootProps();

	const hasError = dropzone.isDragActive && dropzone.isDragReject;
	const isDragAccepted = dropzone.isDragActive && dropzone.isDragAccept;

	return (
		<Button variant="outlined" onClick={onClick}
			sx={{
				cursor: 'pointer',
				width: 1,
				padding: 1,
				borderRadius: 2,
				border: '2px dashed gray',
				display: 'flex',
				alignItems: 'center',
				color: 'text.disabled',
				flexDirection: 'column',
				justifyContent: 'center',

				...(hasError && {
					color: 'error.main',
					bgcolor: 'color-mix(in srgb, var(--mui-palette-error-main), transparent 80%)',
					borderColor: 'var(--mui-palette-error-main)',
				}),

				...(isDragAccepted && {
					color: 'success',
					bgcolor: 'color-mix(in srgb, var(--mui-palette-success-main), transparent 80%)',
					borderColor: 'var(--mui-palette-success-main)',
				}),
				
			}}>
			{props.children}
		</Button>
	);
}