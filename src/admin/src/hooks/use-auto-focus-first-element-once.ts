import { useCallback, useRef } from 'react';

export function useAutoFocusFirstElementOnce() {
	const hasFocusedRef = useRef(false);

	const setRef = useCallback((node: HTMLElement | null) => {
		if (!node || hasFocusedRef.current) return;

		// Skip if something inside the container already has focus
		if (node.contains(document.activeElement)) {
			hasFocusedRef.current = true;
			return;
		}

		const focusableSelectors = [
			'a[href]',
			'button:not([disabled])',
			'input:not([disabled]):not([type="hidden"])',
			'select:not([disabled])',
			'textarea:not([disabled])',
			'[tabindex]:not([tabindex="-1"])'
		];
		const firstFocusable = node.querySelector<HTMLElement>(focusableSelectors.join(', '));
		
		if (firstFocusable) {
			firstFocusable.focus();
			hasFocusedRef.current = true;
		}
	}, []);

	return setRef;
}
