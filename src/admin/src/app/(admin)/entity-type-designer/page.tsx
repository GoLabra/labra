"use client";

import { Box, Container, Stack, Typography } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import CreateIcon from '@mui/icons-material/Create';
import { useDocumentTitle } from "@/hooks/use-document-title";
import { OfficialCommunicationChannels } from "@/shared/components/official-communication-channels";


export default function EntityTypeDesigner() {

    useDocumentTitle({ title: 'Entity Type Designer' });

    return (
        <Box
            sx={{
                flexGrow: 1,
                py: 4
            }}>
            <Container
                maxWidth="md">

                <Box sx={{
                    flexGrow: 1,
                    py: 4,

                }}>
                    <Stack spacing={5}>

                        <Typography
                            color="text.primary"
                            variant="overline">
                            Hello ;)
                        </Typography>

                        <Typography
                            color="text.primary"
                            variant="body2">
                            Once you have created or selected an entity, you can use the various tools provided to manage your content effectively. The management tools are designed to give you full control over your content, ensuring it is accurate and presented in the best possible way.
                        </Typography>

                        <Stack spacing={2}>

                            <Stack
                                alignItems="center"
                                direction="row"
                                spacing={3}
                            >
                                <Typography variant="h3" width="24px">
                                    <AddIcon />
                                </Typography>

                                <div>
                                    <Typography variant="subtitle1" color="text.primary">
                                        Create a New Entity
                                    </Typography>
                                    <Typography
                                        color="text.secondary"
                                        variant="body2"
                                    >
                                        Start by creating a new entity. Click on the &apos;Add Entity&apos; button to begin.
                                    </Typography>
                                </div>
                            </Stack>

                            <Stack
                                alignItems="center"
                                direction="row"
                                spacing={3}>

                                <Typography variant="h3" width="24px">
                                    <CreateIcon fontSize="small"/>
                                </Typography>

                                <div>
                                    <Typography variant="subtitle1" color="text.primary">
                                        Select an Existing Entity
                                    </Typography>

                                    <Typography
                                        color="text.secondary"
                                        variant="body2">
                                        If you have existing entities, you can select one from the sidebar to view and edit its details.
                                    </Typography>
                                </div>
                            </Stack>

                        </Stack>

                        <OfficialCommunicationChannels></OfficialCommunicationChannels>

                    </Stack>

                </Box>
            </Container>
        </Box>
    )
}