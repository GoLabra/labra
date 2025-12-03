import { PageHeader } from "@/shared/components/page-header"
import { Container, Stack, Typography } from "@mui/material"
import { ShowGraphQlQuery } from "./show-graph-ql-query";
import { WithEntitiesSchema } from "./with-entities-schema";

export const DeveloperSchema = () => {

	const EntitiesSchema = WithEntitiesSchema(ShowGraphQlQuery);

	return (
		<Container
			maxWidth="md"
			sx={{
				height: '100%',
				py: 2
			}}>

			<Stack
				spacing={2}
				sx={{ height: '100%' }}>

				<PageHeader
					sx={{
						pl: 1
					}}>
					<Stack
						direction="row"
						justifyContent="space-between"
						alignItems="center"
						spacing={1}>

						<Stack direction="row" alignItems="center" gap={1}>
							<Typography variant="h1">
								GRAPH-QL API
							</Typography>

						</Stack>

						<Stack
							direction="row"
							spacing={1}>

							{/* <Button
                                size="medium"
                                variant="contained"
                                // startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                                onClick={() => exportPostmanCollection()}
                                aria-label="Add new entry"
                                aria-haspopup="dialog">
                                Export as Postman Collection
                            </Button> */}
						</Stack>

					</Stack>
				</PageHeader>

				<Stack spacing={2}>
					<EntitiesSchema />
				</Stack>

			</Stack>

		</Container>
	)
}