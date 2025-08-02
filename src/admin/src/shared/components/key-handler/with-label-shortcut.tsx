import { ComponentType, forwardRef, MutableRefObject, PropsWithChildren, ReactNode } from "react";
import { Key, Modifier } from "./types";
import { useKeyHandler } from "./use-key-handler";
import { ShortcutViewer } from "./shortcut-viewer";
import { Button, ListItem, Stack } from "@mui/material";
import { ActionListItem } from "../action-list-item";

export interface useWithLabelShortcutProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
	onClick: (e: React.KeyboardEvent<HTMLElement> | undefined) => void;
	cancelledShortcutBubble?: boolean;
	target?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
	forceInEditable?: boolean;
	disabled?: boolean;
}
export const useWithLabelShortcut = (props: useWithLabelShortcutProps) => {
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

interface WithLabelShortcutProps {
	shortcutKey: Key;
	shortcutModifiers?: Modifier | Modifier[];
	shortcutTarget?: 'global' | HTMLElement | MutableRefObject<HTMLDivElement | null> | null;
	shortcutForceInEditable?: boolean;
	disabled?: boolean;
	onClick?: (...args: any[]) => void;
	label: ReactNode;
}

export function withLabelShortcut<T>(
	WrappedComponent: ComponentType<T>
) {
	const ComponentWithShortcut = forwardRef<any, T & WithLabelShortcutProps>(
		(props, ref) => {
			const {
				shortcutKey,
				shortcutModifiers,
				shortcutTarget,
				shortcutForceInEditable,
				disabled,
				onClick,
				label,
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
				<WrappedComponent label={<Stack direction="row" alignItems="center">
						{label}
						<ShortcutViewer
							keyToHandle={shortcutKey}
							modifiers={shortcutModifiers}
							sx={{
								marginLeft: '10px',
							}}
						/>
				</Stack>}
				 {...restProps} onClick={onClick} disabled={disabled} ref={ref} />
			);
		}
	);

	ComponentWithShortcut.displayName = `withClickShortcut(${WrappedComponent.displayName || WrappedComponent.name || 'Component'})`;

	return ComponentWithShortcut;
}

export const ShortcutActionItem = withLabelShortcut(ActionListItem);