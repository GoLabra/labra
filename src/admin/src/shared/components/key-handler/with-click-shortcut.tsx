import { ComponentType, forwardRef, MutableRefObject, PropsWithChildren } from "react";
import { Key, Modifier } from "./types";
import { useKeyHandler } from "./use-key-handler";
import { ShortcutViewer } from "./shortcut-viewer";
import { Button, ButtonBase, IconButton, Stack, Tooltip } from "@mui/material";
import { ResponsiveButton } from "@/styles/button.responsive";
import { IconTextButton } from "../IconTextButton";

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
	onClick?: (...args: any[]) => void;
	disabled?: boolean;
	tooltip?: string;
	setRef?: (instance: HTMLElement | null) => void;
}

export function withClickShortcut<T>(
	WrappedComponent: ComponentType<T>,
	showShortcut: boolean = true,
) {
	const ComponentWithShortcut = forwardRef<any, T & PropsWithChildren<WithShortcutProps>>((props, ref) => {
		const {
			shortcutKey,
			shortcutModifiers,
			shortcutTarget,
			shortcutForceInEditable,
			disabled,
			onClick,
			tooltip,
			children,
			setRef,
			...rest
		} = props;

		useKeyHandler({
			keyToHandle: shortcutKey,
			modifiers: shortcutModifiers,
			target: shortcutTarget ?? 'global',
			disabled: disabled,
			forceInEditable: shortcutForceInEditable,
			cancelledShortcutBubble: false,
			onTriggered: () => onClick?.(),
		});

		const restProps = rest as unknown as T;

		const component = (<WrappedComponent {...restProps} onClick={onClick} disabled={disabled} ref={ref}>
				<>
					{children}
					{showShortcut && <ShortcutViewer
						keyToHandle={shortcutKey}
						modifiers={shortcutModifiers}
						sx={{
							marginLeft: '10px',
						}}
					/>}
				</>
			</WrappedComponent>
		);

		return tooltip ? (
			<Tooltip  
				placement='bottom' 
				arrow 
				title={<Stack direction="row" alignItems="center" gap={1}>
							{tooltip}
							<ShortcutViewer 
								keyToHandle={shortcutKey}
								modifiers={shortcutModifiers} />
						</Stack>}>
				{component}
			</Tooltip>
		) : component;
	}
	);

	ComponentWithShortcut.displayName = `withClickShortcut(${WrappedComponent.displayName || WrappedComponent.name || 'Component'})`;

	return ComponentWithShortcut;
}



export const ShortcutButton = withClickShortcut(IconTextButton);
export const ShortcutIconButton = withClickShortcut(IconButton, false);
export const ShortcutButtonBase = withClickShortcut(ButtonBase, false);
