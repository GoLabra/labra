import { Stack, Typography } from "@mui/material"
import { PropsWithChildren } from "react";

interface SectionTitleProps {
	title: string;
}
export const SectionTitle = (props: PropsWithChildren<SectionTitleProps>) => {
	return (
		<Stack direction="row" alignItems="center" mt={3} top="-20px" gap={1} zIndex={10} position="sticky" width={1} bgcolor="background.paper">
			<Typography variant="h3" py={2}>
				{props.title}
			</Typography>
			
			{props.children}

		</Stack>
	)
}