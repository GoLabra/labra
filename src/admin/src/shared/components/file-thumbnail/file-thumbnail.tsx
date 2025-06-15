import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Tooltip from '@mui/material/Tooltip';

import { fileThumbnailClasses } from './classes';
import { fileData, fileThumb, fileFormat } from './utils';
import { RemoveButton, DownloadButton } from './action-buttons';
import { IconButtonProps, SxProps, Theme } from '@mui/material';
import { FileWithPath } from 'react-dropzone';

// ----------------------------------------------------------------------

export type FileWithContent = {
	name: string; 
	content: string;
	size?: number; 
	lastModified?: number;
}

interface FileThumbnailProps {
	sx?: SxProps<Theme>;
	file: FileWithPath | string | FileWithContent;
	tooltip?: boolean;
	onRemove?: () => void;
	imageView?: boolean;
	slotProps?: {
		img?: any;
		icon?: any;
		removeBtn?: IconButtonProps;
		downloadBtn?: IconButtonProps;
	};
	onDownload?: () => void;
}
export function FileThumbnail({
	sx,
	file,
	tooltip,
	onRemove,
	imageView,
	slotProps,
	onDownload,
	...other
}: FileThumbnailProps) {
	//const previewUrl = typeof file === 'string' ? file : URL.createObjectURL(file);

	const { name, path, preview } = fileData(file);

	const format = fileFormat(path || name || '')

	const renderImg = (
		<Box
			component="img"
			src={preview}
			className={fileThumbnailClasses.img}
			sx={{
				width: 1,
				height: 1,
				objectFit: 'cover',
				borderRadius: 'inherit',
				...slotProps?.img,
			}}
		/>
	);

	const renderIcon = (
		<Box
			component="img"
			src={fileThumb(format)}
			className={fileThumbnailClasses.icon}
			sx={{ width: 1, height: 1, ...slotProps?.icon }}
		/>
	);

	const renderContent = (
		<Stack
			component="span"
			className={fileThumbnailClasses.root}
			sx={{
				width: 36,
				height: 36,
				flexShrink: 0,
				borderRadius: 1.25,
				alignItems: 'center',
				position: 'relative',
				display: 'inline-flex',
				justifyContent: 'center',
				...sx,
			}}
			{...other}
		>
			{format === 'image' && imageView ? renderImg : renderIcon}

			{onRemove && (
				<RemoveButton
					onClick={onRemove}
					className={fileThumbnailClasses.removeBtn}
					sx={slotProps?.removeBtn}
				/>
			)}

			{onDownload && (
				<DownloadButton
					onClick={onDownload}
					className={fileThumbnailClasses.downloadBtn}
					sx={slotProps?.downloadBtn}
				/>
			)}
		</Stack>
	);

	if (tooltip) {
		return (
			<Tooltip
				arrow
				title={name}
				slotProps={{ popper: { modifiers: [{ name: 'offset', options: { offset: [0, -12] } }] } }}
			>
				{renderContent}
			</Tooltip>
		);
	}

	return renderContent;
}
