import { Box, BoxProps, styled, SxProps, Theme } from "@mui/material";
import { Key, Modifier } from "./types";
import { useMemo } from "react";
import EnterIcon from "@/assets/icons/labra/enter";

export const ShortcutViewerStyles = styled(Box)<BoxProps>(({ theme }) => ({
	display: 'inline-block',
	padding: '0 5px',
	border: `1px solid color-mix(in srgb, var(--mui-palette-common-onBackground) 10%, #00000000 90%)`,
	minWidth: '20px',
	fontSize: theme.typography.pxToRem(10),
	borderRadius: '2px',
	textAlign: 'center',
	backgroundColor: `color-mix(in srgb, var(--mui-palette-common-onBackground) 5%, #00000000 95%)`,
	textShadow: '0px 1px 1px #0000008c',
	'svg': {
		verticalAlign: 'middle'
	},

    [theme.breakpoints.down('md')]: {
		display: 'none'
    }


})) as typeof Box;


export interface ShortcutViewerProps {
	keyToHandle: Key;
	modifiers?: Modifier | Modifier[];
	sx?: SxProps<Theme>;
}
export const ShortcutViewer = (props: ShortcutViewerProps) => {

	const modifier = useMemo(() => props.modifiers ? Array.isArray(props.modifiers) ? props.modifiers : [props.modifiers] : null, [props.modifiers]);

	const keyToHandle = useMemo(() => {
		if(props.keyToHandle == Key.Enter){
			return <EnterIcon />
		}
		return props.keyToHandle;
	}, [props.keyToHandle]);

	return (
		<ShortcutViewerStyles
			sx={{
				
				...props.sx,
			}}>

			{modifier && (
				modifier?.map(i => i).join(' + ') ?? '')
			}

			{modifier && ' + '}

			{keyToHandle}	

		</ShortcutViewerStyles>
	)
}