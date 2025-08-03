import { useAutoFocusFirstElementOnce } from "@/hooks/use-auto-focus-first-element-once";
import { DialogContent, DialogContentProps } from "@mui/material"
import { PropsWithChildren } from "react"

export const DialogContentWithAutofocus = (props: PropsWithChildren<DialogContentProps>) => {
	const setContainerRef = useAutoFocusFirstElementOnce();
	
	const { children, ...other } = props;
	return (
		<DialogContent
			{...other}
			ref={setContainerRef}
			>
			{props.children}
		</DialogContent>
	)
}