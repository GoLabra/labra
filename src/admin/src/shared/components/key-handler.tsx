import React, { ComponentType, forwardRef, MouseEventHandler, MutableRefObject, PropsWithChildren, Ref, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Box, Button } from '@mui/material';


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
	onTriggered: (e: React.KeyboardEvent<HTMLElement> | undefined) => void;
	cancelledShortcutBubble?: boolean;
	target?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
	forceInEditable?: boolean;
	disabled?: boolean;
}

export const useKeyHandler = (props: KeyHandlerProps) => {

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

		e.shortcutBubbleCancelled = props.cancelledShortcutBubble;

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

export const useWithShortcut = (props: KeyHandlerProps) => {
	const keyHandler = useKeyHandler(props);

	return {
		shortcutKey: props.keyToHandle,
		shortcutModifiers: props.modifiers,
		shortcutTarget: props.target,
		shortcutForceInEditable: props.forceInEditable,
		onClick: (...args: any[]) => {
			props.onTriggered(undefined);
		},
		disabled: props.disabled,
		setRef: keyHandler.setRef
	}
}

interface WithShortcutProps {
  shortcutKey: Key;
  shortcutModifiers?: Modifier | Modifier[];
  shortcutTarget?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
  shortcutForceInEditable?: boolean;
  disabled?: boolean;
  onClick?: (...args: any[]) => void;
}

function withShortcut<T>(
  WrappedComponent: ComponentType<T>
) {
  const ComponentWithShortcut = forwardRef<any, T & PropsWithChildren<WithShortcutProps>>(
    (props, ref) => {
      const {
        shortcutKey,
        shortcutModifiers,
		shortcutTarget,
        shortcutForceInEditable,
		disabled,
        onClick,
        children,
        ...rest
      } = props;

      useKeyHandler({
        keyToHandle: shortcutKey,
        modifiers: shortcutModifiers,
        target: shortcutTarget ?? 'global',
		disabled: disabled,
        forceInEditable: shortcutForceInEditable,
        cancelledShortcutBubble: true,
        onTriggered: () => onClick?.()
      });

      // Cast is safe: rest contains only props for T
      const restProps = rest as unknown as T;

      return (
        <WrappedComponent {...restProps} onClick={onClick} disabled={disabled} ref={ref}>
          <>
            {children}
            <ShortcutView
              keyToHandle={shortcutKey}
              modifiers={shortcutModifiers}
            />
          </>
        </WrappedComponent>
      );
    }
  );

  ComponentWithShortcut.displayName = `withShortcut(${
    WrappedComponent.displayName || WrappedComponent.name || 'Component'
  })`;

  return ComponentWithShortcut;
}

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
				backgroundColor: '#ffffff10',
				textShadow: '0px 1px 2px #0000008c'
			}}>

			{modifier && (
				modifier?.map(i => i).join(' + ') ?? '')
			}

			{modifier && ' + '}

			{props.keyToHandle}

		</Box>
	)
}

export const ShortcutButton = withShortcut(Button);


declare global {
	interface KeyboardEvent {
		shortcutBubbleCancelled?: boolean;
	}
}