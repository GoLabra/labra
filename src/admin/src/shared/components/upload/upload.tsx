import { forwardRef, memo, useEffect, useMemo, useState } from 'react';
import { DropzoneState, useDropzone } from 'react-dropzone';
import { Box, Typography } from '@mui/material';
import { FileDrop } from './components/file-drop';
import { FilesPreview } from './components/files-preview';
import { BrowseButton } from './components/browse-button';
import AttachFileIcon from '@mui/icons-material/AttachFile';

interface UploadFilesProps {
	name: string;
	dropzone: DropzoneState;
	value: (File | string)[];
	viewMode?: 'list' | 'grid';
	maxFiles?: number;
	onRemove?: (file: File | string) => void;
}
export const UploadFiles = (props: UploadFilesProps) => {

	const { dropzone } = props;

	const isFull = props.maxFiles != undefined && props.value?.length >= props.maxFiles;

	return (
		<FileDrop name={props.name} dropzone={dropzone}>
			<FilesPreview
				viewMode={props.viewMode}
				onRemove={props.onRemove}
				firstNode={isFull == false ? <BrowseButton dropzone={dropzone}>
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
				</BrowseButton> : undefined}

				files={props.value}
			/>
		</FileDrop>
	);
};