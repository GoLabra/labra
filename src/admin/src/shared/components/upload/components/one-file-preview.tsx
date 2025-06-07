import { memo } from "react";
import { isImageFile } from "../../is-image-file";
import { Box, Stack } from "@mui/material";
import { Image } from "@/shared/components/image/image";
import { FileThumbnail } from "../../file-thumbnail/file-thumbnail";

interface FilePreviewProps {
	file?: File | string;
}
export const OneFilePreview = memo((props: FilePreviewProps) => {  

	if (!props.file) {
		return null;
	}

	const isImage = isImageFile(props.file);

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
});