export const isJsonOrNull = (val?: string | null) => {
	if (!val) {
		return true;
	}
	try {
		JSON.parse(val);
		return true;
	} catch {
		return false;
	}
}

export const toJsonOrNull = (val?: string | null) => {
	if (!val) {
		return null
	}
	try {
		return JSON.parse(val);
	} catch {
		return null;
	}
}