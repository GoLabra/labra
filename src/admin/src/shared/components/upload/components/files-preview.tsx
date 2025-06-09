import Box from '@mui/material/Box';
import IconButton from '@mui/material/IconButton';
import ListItemText from '@mui/material/ListItemText';
// import { varAlpha } from 'src/theme/styles'; 
import { fileData, FileThumbnail } from '../../file-thumbnail';
import { fData } from '@/lib/utils/format-number';
import { Stack, SxProps, Theme } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import CancelIcon from '@mui/icons-material/Cancel';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';

interface FilePreviewProps {
	sx?: SxProps<Theme>;
	onRemove?: (file: File | string) => void;
	lastNode?: React.ReactNode;
	viewMode?: 'list' | 'grid';
	slotProps?: {
		thumbnail?: any;
	};
	firstNode?: React.ReactNode;
	files?: any[];
}
export function FilesPreview(props: FilePreviewProps) {

	const viewMode = props.viewMode ?? 'grid';

	const renderFirstNode = props.firstNode && (
		<Box
			component="li"
			sx={{
				listStyleType: 'none',
				...(viewMode === 'grid' && {
					width: 100,
					height: 100,
					display: 'inline-flex',
				}),
			}}
		>
			{props.firstNode}
		</Box>
	);

	const renderLastNode = props.lastNode && (
		<Box
			component="li"
			sx={{
				listStyleType: 'none',
				...(viewMode === 'grid' && {
					width: 100,
					height: 100,
					display: 'inline-flex'
				}),
			}}
		>
			{props.lastNode}
		</Box>
	);

	return (
		<Box
			component="ul"
			sx={{
				gap: 1,
				display: 'flex',
				padding: 0,
				margin: 0,
				flexDirection: 'column',
				...(viewMode === 'grid' && {
					flexWrap: 'wrap',
					flexDirection: 'row',
				}),
				...props.sx,
			}}
		>
			{renderFirstNode}

			{props.files?.map((file) => {
				const { name, size } = fileData(file);

				if (viewMode === 'grid') {
					return (
						<Box component="li" key={name} sx={{
							width: 100,
							height: 100,
							display: 'inline-flex',
							border: (theme) => `solid 1px ${theme.palette.grey[500]}`,
							position: 'relative',

							borderRadius: 1,
							borderStyle: 'solid',
							borderWidth: '1px',
							borderColor: 'color-mix(in srgb, var(--mui-palette-grey-500), transparent 80%)',
							padding: '1px'
						}}>

							{props.onRemove && (<Box sx={{position: 'absolute', top: 0, right: 0, zIndex: 1, padding: '1px'}} onClick={() => props.onRemove?.(file)}>
								<IconButton aria-label="delete" size="small">
									<HighlightOffIcon fontSize="inherit" />
								</IconButton>
							</Box>)}

							<FileThumbnail
								tooltip
								imageView={true}
								file={file}
								onRemove={() => props.onRemove?.(file)}
								sx={{
									width: 1,
									height: 1,
									//   border: (theme) =>
									//     `solid 1px ${theme.palette.grey[500]}`,
								}}
								slotProps={{ icon: { width: 36, height: 36 } }}
								{...props.slotProps?.thumbnail}
							/>
						</Box>
					);
				}

				return (
					<Stack
						component="li"
						key={name}
						direction="row"
						gap={1}
						alignItems="center"
					>
						<FileThumbnail imageView file={file} {...props.slotProps?.thumbnail} />

						<ListItemText
							primary={name}
							secondary={fData(size)}
							// secondaryTypographyProps={{ component: 'span', typography: 'caption' }}

							slotProps={{
								primary: {
									typography: {
										textOverflow: 'ellipsis',
										fontSize: '0.9rem',
									},
								},
								secondary: {
									component: 'span',
									typography: 'caption',
								}

							}}

						/>

						{props.onRemove && (
							<IconButton size="small" onClick={() => props.onRemove?.(file)}>
								<HighlightOffIcon fontSize="inherit" />
							</IconButton>
						)}
					</Stack>
				);
			})}

			{renderLastNode}
		</Box>
	);
}
