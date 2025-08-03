import { useAutoFocusFirstElementOnce } from "@/hooks/use-auto-focus-first-element-once";
import { DialogContent, DialogContentProps } from "@mui/material"
import { PropsWithChildren } from "react"

export const DialogContentWithAutofocus = (props: PropsWithChildren<DialogContentProps>) => {
	const autoFocusHandler = useAutoFocusFirstElementOnce();
	
	const { children, ...other } = props;
	return (
		<DialogContent
			{...other}
			ref={autoFocusHandler.setRef}
			>
			{props.children}
		</DialogContent>
	)
}