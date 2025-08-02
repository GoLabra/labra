
export enum Key {
	Enter = 'Enter',
	Escape = 'Escape',
	Space = ' ',
	ArrowUp = 'ArrowUp',
	ArrowDown = 'ArrowDown',
	ArrowLeft = 'ArrowLeft',
	ArrowRight = 'ArrowRight',
	n = 'n',
	s = 's',
	slash = '/',
	// Add more keys as needed
}

export type Modifier = "Ctrl" | "Shift" | "Alt" | "Meta";

declare global {
	interface KeyboardEvent {
		shortcutBubbleCancelled?: boolean;
	}
}