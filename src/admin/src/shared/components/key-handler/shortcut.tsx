import { ShortcutViewer, ShortcutViewerProps } from "./shortcut-viewer";
import { useKeyHandler, useKeyHandlerProps } from "./use-key-handler";

interface ShortcutProps extends useKeyHandlerProps
{ 
	slotProps?: {
		ShortcutViewer?: Partial<Omit<ShortcutViewerProps, 'keyToHandle'>>
	}
}

export const Shortcut = (props: ShortcutProps) => {

	useKeyHandler({...props, 
		target: props.target ?? 'global'
	});
	
	return (<ShortcutViewer 
				keyToHandle={props.keyToHandle} 
				modifiers={props.modifiers}
				{...props.slotProps?.ShortcutViewer} />);

}