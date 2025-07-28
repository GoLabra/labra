import { FileWithPath } from "react-dropzone";
import { FileData } from "./file-thumbnail";

// Define more types here
const FORMAT_PDF = ['pdf'];
const FORMAT_TEXT = ['txt'];
const FORMAT_PHOTOSHOP = ['psd'];
const FORMAT_WORD = ['doc', 'docx'];
const FORMAT_EXCEL = ['xls', 'xlsx'];
const FORMAT_ZIP = ['zip', 'rar', 'iso'];
const FORMAT_ILLUSTRATOR = ['ai', 'esp'];
const FORMAT_POWERPOINT = ['ppt', 'pptx'];
const FORMAT_AUDIO = ['wav', 'aif', 'mp3', 'aac'];
const FORMAT_IMG = ['jpg', 'jpeg', 'gif', 'bmp', 'png', 'svg', 'webp'];
const FORMAT_VIDEO = ['m4v', 'avi', 'mpg', 'mp4', 'webm'];

const iconUrl = (icon: string) => `/assets/icons/files/${icon}.svg`;

export function fileIcon(extension?: string | null) {
	
	if(!extension){
		return iconUrl('ic-file');
	}

	switch (true) {
		case extension === 'folder':
		 	return iconUrl('ic-folder');

		case FORMAT_TEXT.includes(extension):
			return iconUrl('ic-txt');

		case FORMAT_AUDIO.includes(extension):
			return iconUrl('ic-audio');

		case FORMAT_ZIP.includes(extension):
			return iconUrl('ic-zip');

		case FORMAT_IMG.includes(extension):
			return iconUrl('ic-img');

		case FORMAT_VIDEO.includes(extension):
			return iconUrl('ic-video');

		case FORMAT_WORD.includes(extension):
			return iconUrl('ic-word');

		case FORMAT_EXCEL.includes(extension):
			return iconUrl('ic-excel');

		case FORMAT_POWERPOINT.includes(extension):
			return iconUrl('ic-power_point');

		case FORMAT_PDF.includes(extension):
			return iconUrl('ic-pdf');

		case FORMAT_PHOTOSHOP.includes(extension):
			return iconUrl('ic-pts');

		case FORMAT_ILLUSTRATOR.includes(extension):
			return iconUrl('ic-ai');

		default:
			return iconUrl('ic-file');
	}
}

// ----------------------------------------------------------------------

export function fileIsImage(extension?: string | null) {
	if(!extension){
		return false;
	}
	return FORMAT_IMG.includes(extension);
}

export function fileTypeByUrl(fileUrl?: string | null): string | null {
	if(!fileUrl){
		return null;
	}
	return (fileUrl && fileUrl.split('.').pop())?.toLowerCase() || '';
}

// ----------------------------------------------------------------------

export function fileNameByUrl(fileUrl: string) {
	return fileUrl.split('/').pop();
}

// ----------------------------------------------------------------------

export function fileData(file: FileWithPath | FileData | string): FileData {

	// From content
	if (file && typeof file === 'object' && !('path' in file)) { 
		return file;
	}
	
	// From url
	if (typeof file === 'string') {
		return {
			preview: file,
			name: fileNameByUrl(file),
			size: undefined,
			lastModified: undefined,
		};
	}


	const isImage = fileIsImage(fileTypeByUrl(file.name));
	// From file
	return {
		name: file.name,
		size: file.size,
		preview: isImage ? URL.createObjectURL(file) : undefined,
		lastModified: file.lastModified,
	};
}
