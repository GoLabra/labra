import { MutableRefObject, useCallback, useEffect, useRef, useState } from "react";
import { Key, Modifier } from "./types";

export interface useKeyHandlerProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
	onTriggered: (e: React.KeyboardEvent<HTMLElement> | undefined) => void;
	cancelledShortcutBubble?: boolean;
	target?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
	forceInEditable?: boolean;
	disabled?: boolean;
}

export const useKeyHandler = (props: useKeyHandlerProps) => {

	// Refs to hold the latest values
	const propRef = useRef(props);
	const [manualTarget, setManualTarget] = useState<HTMLElement | null>(null);

	// Update refs when props change
	useEffect(() => {
		propRef.current = props;
	}, [props]);

	const handleKeyDown = useCallback((e: KeyboardEvent) => {

		if (propRef.current.disabled) {
			return;
		}

		if (e.shortcutBubbleCancelled) {
			return;
		}

		e.shortcutBubbleCancelled = propRef.current.cancelledShortcutBubble;

		const target = e.target as HTMLElement;
		const isEditable =
			target.tagName === 'INPUT' ||
			target.tagName === 'TEXTAREA' ||
			target.isContentEditable;

		if (isEditable && propRef.current.forceInEditable != true) return;

		const modifier = propRef.current.modifiers ? Array.isArray(propRef.current.modifiers) ? propRef.current.modifiers : [propRef.current.modifiers] : null;

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


		if (modifiersPressed && e.key === propRef.current.keyToHandle) {
			propRef.current.onTriggered(e as unknown as React.KeyboardEvent<HTMLDivElement>);
			e.preventDefault();
			e.shortcutBubbleCancelled = true;
		}
	}, []);

	useEffect(() => {

		const target = (() => {

			if (manualTarget) {
				return manualTarget;
			}

			const propsTarget = propRef.current.target;
			if (!propsTarget) {
				return;
			}

			if (propsTarget == 'global') {
				return document.body;
			}

			if (propsTarget instanceof HTMLElement) {
				return propsTarget;
			}

			if ('current' in propsTarget) {
				return propsTarget.current;
			}

			return null;
		})();

		target?.addEventListener('keydown', handleKeyDown);

		return () => {
			target?.removeEventListener('keydown', handleKeyDown);
		}

	}, [handleKeyDown, props.target, manualTarget]);


	return {
		setRef: (instance: HTMLElement | null) => {
			setManualTarget(instance);
		}
	};
};