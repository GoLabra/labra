'use client'

import { AuthContextType } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { dynamicLayoutItem } from "@/core-features/dynamic-layout/src/dynamic-layout";
import { paths } from "@/lib/paths";
import { Copywrite } from "@/shared/components/copywrite";
import { Alert, Box, Button, Stack, TextField, Typography } from "@mui/material";
import NextLink from 'next/link';
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, useCallback } from "react";
import { PRODUCT_NAME } from '@/config/CONST';
import { gql } from "@apollo/client";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { PasswordFormField } from "@/core-features/dynamic-form/form-fields/PasswordField";


const schema = z.object({
    email: z.string().min(1, 'Email is required').email(),
    firstName: z.string().min(1, 'First Name is required'),
    lastName: z.string().min(1, 'Last Name is required'),
    password: z.string().min(1, 'Password is required'),
});

export default function AuthPage() {

    const auth = useAuth<AuthContextType>();
    const router = useRouter();
    const searchParams = useSearchParams();
    const returnTo = searchParams.get('returnTo') || undefined;

    const formMethods = useForm({
        resolver: zodResolver(schema)
    });

    const onSubmit = useCallback(async (data: any) => {
        await auth.signUp(data);
        router.push(returnTo || paths.auth.jwt.login);
    }, [auth.signUp]);

    return (<>

        <Box sx={{ pt: 7, pb: 14 }}>

            <Typography textAlign="center" fontWeight={100} fontSize={48} fontFamily="Roboto" sx={{ mb: 7 }}>
                {PRODUCT_NAME}
            </Typography>

            <Stack
                alignItems="center"
                direction="row"
                justifyContent="space-between"
                spacing={1}
                sx={{ mb: 3 }}
            >
                <Typography variant="h4">
                    Superuser Sign Up
                </Typography>
                <Button
                    component={NextLink}
                    href={paths.auth.jwt.login}
                >
                    Sign In
                </Button>
            </Stack>

            <Form methods={formMethods} onSubmit={formMethods.handleSubmit(onSubmit)} >
                <Stack gap={1.5}>
                    <TextShortFormField name="email" label="Email" required />
                    <TextShortFormField name="firstName" label="First Name" required />
                    <TextShortFormField name="lastName" label="Last Name" required />
                    <PasswordFormField name="password" label="Password" required />
                </Stack>

                <Button
                    fullWidth
                    sx={{ mt: 3 }}
                    type="submit"
                    variant="contained"
                >
                    Sign Up
                </Button>

            </Form>

            <Box textAlign="center" sx={{ mt: 3 }}>
                <Copywrite />
            </Box>

        </Box>
    </>
    )
}
