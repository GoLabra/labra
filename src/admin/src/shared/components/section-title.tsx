import { Stack, Typography } from "@mui/material"
import { PropsWithChildren } from "react";

interface SectionTitleProps {
	title: string;
	endAdornment?: React.ReactNode;
}
export const SectionTitle = (props: PropsWithChildren<SectionTitleProps>) => {
	return (
		<Stack direction="row" alignItems="center" mt={3} top="-20px" justifyContent="space-between" gap={1} zIndex={10} position="sticky" width={1} bgcolor="background.paper">

			<Stack direction="row" alignItems="center" gap={1}>
				<Typography variant="h3" py={2}>
					{props.title}
				</Typography>
				{props.children}
			</Stack>
			
			{props.endAdornment}

		</Stack>
	)
}