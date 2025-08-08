import { styled, Chip, ComponentsProps, ComponentsOverrides, ComponentsVariants, ChipProps, Box, Stack, BoxProps } from "@mui/material";
import { forwardRef, PropsWithChildren, ReactNode, useMemo } from "react";
import { useTheme } from "@mui/material/styles";


const CounterGroupRoot = styled(Stack, {
	name: 'CounterGroup',
	slot: 'root'
})({
	'& .Counter-root:not(:first-child):not(:last-child)': {
    borderRadius: 0, // Reset for middle children
  },
  
  '& .Counter-root:first-child:not(:last-child)': {
    borderTopRightRadius: 0,
    borderBottomRightRadius: 0,
  },
  
  '& .Counter-root:last-child:not(:first-child)': {
    borderTopLeftRadius: 0,
    borderBottomLeftRadius: 0,
  }
});


export interface CounterGroupProps extends PropsWithChildren {
	size?: 'small' | 'medium' | 'large' | string & {};
}
const sizeStyles = {
	small: { fontSize: '12px' },
	medium: { fontSize: '14px' },
	large: { fontSize: '16px' }
} as const;

export const CounterGroup = forwardRef<HTMLDivElement, CounterGroupProps>((props, ref) => {
	const { size = 'medium', children } = props;
	
	const sizeStyleProps = useMemo(() => {
		if (size === 'small' || size === 'medium' || size === 'large') {
			return sizeStyles[size as 'small' | 'medium' | 'large'];
		}
		
		if (typeof size === 'string') {
			return {
				fontSize: size
			};
		}
		
		return sizeStyles.medium;
	}, [size]);
	
	return (
		<CounterGroupRoot
			ref={ref}
			direction="row"
			divider={<Box sx={{ width: '1px', backgroundColor: 'var(--mui-palette-divider)' }} />}
			sx={sizeStyleProps}
		>
			{children}
		</CounterGroupRoot>
	);
});



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
export interface CounterProps extends BoxProps {
	label: ReactNode;
}
export const Counter = forwardRef<HTMLDivElement, CounterProps>(({ label, className, ...rest }, ref) => (
	<CounterRoot
		ref={ref}
		{...rest}
		className={`Counter-root ${className ?? ''}`}
	>
		{label}
	</CounterRoot>
));