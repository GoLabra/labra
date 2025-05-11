'use client'

import { useDocumentTitle } from "@/hooks/use-document-title";
import { Container, Stack, Box, Card, CardContent, Typography, Grid, Paper, styled, Button, SvgIcon, Tabs, Tab } from "@mui/material";
import { PageHeader } from "@/shared/components/page-header";
import { UserAccount } from "@/features/user-account";
import { useTabs } from "@/hooks/useTabs";
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import SecurityIcon from '@mui/icons-material/Security';
import { AccountGeneralTab } from "@/features/user-account/account-general-tab";
import { useForm } from "react-hook-form";
import { Form, FormDynamicContextType, SchemaFieldOptions, useSchemaOptions } from "@/core-features/dynamic-form2/dynamic-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { ChangePasswordTab } from "@/features/user-account/change-password-tab";


export default function ContentManager() {

    useDocumentTitle({ title: 'User Account' });
    const tabs = useTabs('general');


    return (

        <Container
            maxWidth="md"
            sx={{
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
                        align-items="flex-start"
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                        spacing={1}>

                        <Typography variant="h2">
                            User Account
                        </Typography>

                        <Stack
                            direction="row"
                            spacing={1}>

                            <Tabs {...tabs} sx={{ mb: { xs: 3, md: 5 } }}>
                                <Tab label="General" icon={<AssignmentIndIcon />} value="general" />
                                <Tab label="Security" icon={<SecurityIcon />} value="security" />
                            </Tabs>
                        </Stack>

                    </Stack>
                </PageHeader>

                {tabs.value === 'general' && <AccountGeneralTab />}
                {tabs.value === 'security' && <ChangePasswordTab />}


            </Stack>
        </Container>

    )
}