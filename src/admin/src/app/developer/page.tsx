'use client'

import { useDocumentTitle } from "@/hooks/use-document-title";
import { Container, Stack, Box, Card, CardContent, Typography, styled, Paper } from "@mui/material";
import { PageHeader } from "@/shared/components/page-header";

const BulletItem = styled(Stack)(({ theme }) => ({
    border: '2px dashed #80808024',
    borderRadius: '3px',
    padding: '10px',
    zIndex: 1,
    backgroundColor: 'var(--mui-palette-background-default)',
}));

export default function ContentManager() {

    useDocumentTitle({ title: 'Developer' });

    return (
        <>
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
                                    Welcome to the Developer Section
                                </Typography>
                            </Stack>

                            <Stack
                                direction="row"
                                spacing={1}>

                                {/* <Button
                                        size="medium"
                                        variant="contained"
                                        startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                                        onClick={() => addNewEntry()}
                                        aria-label="Add new entry"
                                        aria-haspopup="dialog">
                                        Add New
                                    </Button> */}
                            </Stack>

                        </Stack>
                    </PageHeader>

                    <Card>
                        <CardContent>
                            <Typography
                                color="text.primary"
                                variant="body2">
                                Your toolkit for seamless API development. Access the GraphQL playground to test queries, explore responses, and fine-tune your requests. Export collections to Postman, analyze data, and streamline your workflow—all in one place.
                            </Typography>
                        </CardContent>
                    </Card>

                </Stack>
            </Container>
        </>
    )
}

