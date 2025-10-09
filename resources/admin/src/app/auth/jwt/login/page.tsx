'use client'

import { AuthContextType } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { paths } from "@/lib/paths";
import { Copywrite } from "@/shared/components/copywrite";
import { Alert, Box, Button, Stack, TextField, Typography } from "@mui/material";
import NextLink from 'next/link';
import { useRouter, useSearchParams } from "next/navigation";
import { FormEvent, useCallback } from "react";
import { PRODUCT_NAME } from '@/config/CONST';
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { PasswordFormField } from "@/core-features/dynamic-form/form-fields/PasswordField";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";


const schema = z.object({
    email: z.string().min(1, 'Email is required').email(),
    password: z.string().min(1, 'Password is required'),
});

export default function AuthPage() {

    const auth = useAuth<AuthContextType>();
    const router = useRouter();
    const searchParams = useSearchParams();
    const returnTo = searchParams.get('returnTo') || undefined;

    const onSubmit = useCallback(async (data: any) => {
        await auth.signIn(data);
        router.push(returnTo || paths.contentManager.index);
    }, [auth.signIn]);

    const formMethods = useForm({
        resolver: zodResolver(schema)
    });

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
                    Login
                </Typography>
                <Button
                        component={NextLink}
                        href={paths.auth.signup}
                    >
                        Sign Up
                    </Button>
            </Stack>

            <Form methods={formMethods} onSubmit={formMethods.handleSubmit(onSubmit)} >
                <Stack gap={1.5}>
                    <TextShortFormField name="email" label="Email" required />
                    <PasswordFormField name="password" label="Password" required />
                </Stack>

                <Button
                    fullWidth
                    sx={{ mt: 3 }}
                    type="submit"
                    variant="contained"
					disabled={formMethods.formState.isSubmitting}
                >
                    Login
                </Button>

            </Form>

            <Box textAlign="center" sx={{ mt: 3 }}>
                <Copywrite />
            </Box>

        </Box>
    </>
    )
}
