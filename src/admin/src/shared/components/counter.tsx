import { styled, Chip, ComponentsProps, ComponentsOverrides, ComponentsVariants, ChipProps, Box } from "@mui/material";
import { forwardRef, ReactNode } from "react";
import { useTheme } from "@mui/material/styles";

// Create the styled component with proper theme integration
const CounterRoot = styled(Box, {
	name: 'Counter',
	slot: 'root'
})({
	backgroundColor: 'var(--mui-palette-action-disabledBackground)',
	padding: '0px 10px',
    borderRadius: '16px',
	fontWeight: 100
});

// Extend the props interface
export interface CounterProps {
	label: ReactNode;
}
export const Counter = forwardRef<HTMLDivElement, CounterProps>((props, ref) => (	
	<CounterRoot
		ref={ref}
	>
		{props.label}
	</CounterRoot>
));