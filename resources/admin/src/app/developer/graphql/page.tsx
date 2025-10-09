'use client'

import { useDocumentTitle } from "@/hooks/use-document-title";
import { Container, Stack, Box, Card, CardContent, Typography, styled, Button } from "@mui/material";
import { PageHeader } from "@/shared/components/page-header";
import { WithEntitySchema } from "@/features/developer/with-entity-schema";
import { ShowGraphQlQuery } from "@/features/developer/show-graph-ql-query";
import { WithEntitiesSchema } from "@/features/developer/with-entities-schema";

export default function DeveloperGraphQl() {

    useDocumentTitle({ title: 'Developer - Export' });

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

