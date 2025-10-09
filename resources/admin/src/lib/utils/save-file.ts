export function saveFile(filename: string, input: File | string) {
	let blob;

	if (input instanceof File) {
		blob = input;
		filename = input.name || filename;
	} else if (typeof input === 'string') {
		if (input.startsWith('data:')) {
			const [meta, base64Data] = input.split(',');
			const byteString = atob(base64Data);
			const byteArray = new Uint8Array(byteString.length);
			for (let i = 0; i < byteString.length; i++) {
				byteArray[i] = byteString.charCodeAt(i);
			}
			blob = new Blob([byteArray]);
		} else {
			blob = new Blob([input]);
		}
	} else {
		console.error('Input must be a File or string');
		return;
	}

	const link = document.createElement('a');
	link.href = URL.createObjectURL(blob);
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(link.href);
}