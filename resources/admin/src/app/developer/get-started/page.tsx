'use client'

import { useDocumentTitle } from "@/hooks/use-document-title";
import { Container, Stack, Box, Card, CardContent, Typography, styled } from "@mui/material";
import FiberManualRecordIcon from '@mui/icons-material/FiberManualRecord';
import { useEffect, useState } from "react";
import { OfficialCommunicationChannels } from "@/shared/components/official-communication-channels";

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
                            variant="body2">
                            Manage your content effortlessly by selecting an entity from the sidebar. This page is designed to help you oversee, organize, and maintain all your entries effectively.
                        </Typography>

                     
                        <OfficialCommunicationChannels></OfficialCommunicationChannels>
                    </Stack>

                </Box>
            </Container>
        </Box>
    )
}

