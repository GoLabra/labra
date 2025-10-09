export function isImageFile(file?: File | string): boolean {
	if(!file){
		return false;
	}

	if (file instanceof File) {
		return file.type.startsWith('image/');
	} 
	
	if (typeof file === 'string') {
		// For data URLs
		if (file.startsWith('data:')) {
			return file.startsWith('data:image/');
		}
		// For regular URLs, you might want to check the extension
		return /\.(jpg|jpeg|png|gif|webp|bmp)$/i.test(file);
	}
	
	return false;
}