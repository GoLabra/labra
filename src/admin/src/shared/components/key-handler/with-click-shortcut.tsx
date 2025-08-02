import { ComponentType, forwardRef, MutableRefObject, PropsWithChildren } from "react";
import { Key, Modifier } from "./types";
import { useKeyHandler } from "./use-key-handler";
import { ShortcutViewer } from "./shortcut-viewer";
import { Button } from "@mui/material";
import { ResponsiveButton } from "@/styles/button.responsive";

export interface useWithClickShortcutProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
	onClick: (e: React.KeyboardEvent<HTMLElement> | undefined) => void;
	cancelledShortcutBubble?: boolean;
	target?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
	forceInEditable?: boolean;
	disabled?: boolean;
}
export const useWithClickShortcut = (props: useWithClickShortcutProps) => {
	const { onClick, ...rest } = props;
	const keyHandler = useKeyHandler({ ...rest, onTriggered: onClick });

	return {
		shortcutKey: props.keyToHandle,
		shortcutModifiers: props.modifiers,
		shortcutTarget: props.target,
		shortcutForceInEditable: props.forceInEditable,
		onClick: (...args: any[]) => {
			props.onClick(undefined);
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

export function withClickShortcut<T>(
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
				cancelledShortcutBubble: false,
				onTriggered: () => onClick?.()
			});

			const restProps = rest as unknown as T;

			return (
				<WrappedComponent {...restProps} onClick={onClick} disabled={disabled} ref={ref}>
					<>
						{children}
						<ShortcutViewer
							keyToHandle={shortcutKey}
							modifiers={shortcutModifiers}
							sx={{
								marginLeft: '10px',
							}}
						/>
					</>
				</WrappedComponent>
			);
		}
	);

	ComponentWithShortcut.displayName = `withClickShortcut(${WrappedComponent.displayName || WrappedComponent.name || 'Component'})`;

	return ComponentWithShortcut;
}



export const ShortcutButton = withClickShortcut(Button);



