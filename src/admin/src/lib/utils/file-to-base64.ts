export const fileToBase64 = (file?: File) => new Promise((resolve, reject) => { 
	if (!file) {
		return null;
	}
	const reader = new FileReader();
	reader.readAsDataURL(file);
	reader.onload = () => resolve(reader.result);
	reader.onerror = reject;
	
	return reader;
});


export const fileToBase64Sync = (file?: File): string | null => {  
	if (!file) {
		return null;
	}
	const reader = new FileReader();
	reader.readAsDataURL(file);

	return reader.result as string;
};