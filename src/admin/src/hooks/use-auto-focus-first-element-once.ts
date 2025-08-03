import { useCallback, useMemo, useRef } from 'react';

export function useAutoFocusFirstElementOnce() {
	const hasFocusedRef = useRef(false);
	const containerRef = useRef<HTMLElement | null>(null);

	const setRef = useCallback((node: HTMLElement | null) => {
		containerRef.current = node;

		if (!hasFocusedRef.current) {
			focus();
		}
	}, []);

	const focus = useCallback(() => {
		const node = containerRef.current;

		if (!node) {
			return;
		}

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

	return useMemo(() => ({
		setRef,
		focus: () => {
			setTimeout(() => {
				focus();
			}, 0);
		}
	}), [setRef, focus]);

}
