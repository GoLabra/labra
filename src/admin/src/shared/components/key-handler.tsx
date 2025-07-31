import React, { PropsWithChildren, useCallback, useEffect, useMemo, useRef } from 'react';
import { Box } from '@mui/material';


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
	// Add more keys as needed
}

export type Modifier = "Ctrl" | "Shift" | "Alt" | "Meta"; 

interface KeyHandlerProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
	onTriggered: (e: React.KeyboardEvent<HTMLElement>) => void;
	cancelledShortcutBubble?: boolean;
	global?: boolean;
	forceInEditable?: boolean;
	disabled?: boolean;
}

export const useKeyHandler = (props: KeyHandlerProps) => {
	const { keyToHandle, modifiers, onTriggered, disabled, forceInEditable } = props;
	const containerRef = useRef<HTMLElement | null>();

	// Refs to hold the latest values
	const keyToHandleRef = useRef<Key>(keyToHandle);
	const modifiersRef = useRef<Modifier | Modifier[] | undefined>(modifiers);
	const onTriggeredRef = useRef<(e: React.KeyboardEvent<HTMLElement>) => void>(onTriggered);
	const forceInEditableRef = useRef<boolean|undefined>(props.forceInEditable);
	const disabledRef = useRef<boolean|undefined>(props.disabled);

	// Update refs when props change
	useEffect(() => {
		keyToHandleRef.current = keyToHandle;
		modifiersRef.current = modifiers;
		onTriggeredRef.current = onTriggered;
		forceInEditableRef.current = forceInEditable;
		disabledRef.current = disabled;
	}, [keyToHandle, modifiers, onTriggered, forceInEditable, disabled]);

	const handleKeyDown = useCallback((e: KeyboardEvent) => {

		if(disabledRef.current){
			return;
		}

		if (e.shortcutBubbleCancelled) {
			return;
		}

		e.shortcutBubbleCancelled = props.cancelledShortcutBubble;

		const target = e.target as HTMLElement;
		const isEditable = 
			target.tagName === 'INPUT' ||
			target.tagName === 'TEXTAREA' ||
			target.isContentEditable;

		if (isEditable && forceInEditableRef.current != true) return;

		const modifier = modifiersRef.current ? Array.isArray(modifiersRef.current) ? modifiersRef.current : [modifiersRef.current] : null;

		const modifiersPressed = !modifier || modifier.every(modifier => {
			switch (modifier.toLowerCase()) {
				case 'ctrl':
					return e.ctrlKey;
				case 'shift':
					return e.shiftKey;
				case 'alt':
					return e.altKey;
				case 'meta':
					return e.metaKey; // For Command key on Mac or Windows key
				default:
					return false;
			}
		});


		if (modifiersPressed && e.key === keyToHandleRef.current) {
			//e.preventDefault();

			if (onTriggeredRef.current) {
				onTriggeredRef.current(e as unknown as React.KeyboardEvent<HTMLDivElement>);
			}
		}
	}, []); // Empty dependency array ensures the function reference stays the same

	const setRef = useCallback((node: HTMLElement | null) => {
		if (containerRef.current) {
			containerRef.current.removeEventListener('keydown', handleKeyDown);
		}

		containerRef.current = node;

		if (containerRef.current) {
			containerRef.current.addEventListener('keydown', handleKeyDown);
		}
	}, [handleKeyDown]);

	useEffect(() => {
		if (props.global) {
			setRef(document.body);
		}
	}, []);

	return {
		setRef
	};
};


interface ShortcutViewProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
}
export const ShortcutView = (props: ShortcutViewProps) => {

	const modifier = useMemo(() => props.modifiers ? Array.isArray(props.modifiers) ? props.modifiers : [props.modifiers] : null, [props.modifiers]); 

	return (
		<Box
			sx={{
				display: 'inline-block',
				border: '1px solid white',
				padding: '0 5px',
				marginLeft: '10px',
				minWidth: '20px',
				fontSize: '10px',
				borderRadius: '2px',
				textAlign: 'center',
			}}>


			{ modifier && (
				modifier?.map(i => i).join(' + ') ?? '') 
			}
			
			{ modifier && ' + ' }

			{props.keyToHandle}
			
		</Box>
	)
}


declare global {
	interface KeyboardEvent {
		shortcutBubbleCancelled?: boolean;
	}
}