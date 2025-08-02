
export enum Key {
	Enter = 'Enter',
	Escape = 'Escape',
	Space = ' ',
	ArrowUp = 'ArrowUp',
	ArrowDown = 'ArrowDown',
	ArrowLeft = 'ArrowLeft',
	ArrowRight = 'ArrowRight',
	a = 'a',
	c = 'c',
	f = 'f',
	n = 'n',
	o = 'o',
	r = 'r',
	s = 's',
	u = 'u',
	slash = '/',
	// Add more keys as needed
}

export type Modifier = "Ctrl" | "Shift" | "Alt" | "Meta";

declare global {
	interface KeyboardEvent {
		shortcutBubbleCancelled?: boolean;
	}
}