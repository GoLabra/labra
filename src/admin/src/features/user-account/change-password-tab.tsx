import { AuthContextType } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { PasswordFormField } from "@/core-features/dynamic-form/form-fields/PasswordField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Avatar, Box, Button, Card, CardContent, Grid, Stack, Typography } from "@mui/material"
import { useForm } from "react-hook-form";
import { z } from "zod";


const schema = z
    .object({
        oldPassword: z
            .string()
            .min(1, { message: 'Password is required!' }),
        newPassword: z.string().min(1, { message: 'New password is required!' }),
        confirmNewPassword: z.string().min(1, { message: 'Confirm password is required!' }),
    })
    .refine((data) => data.oldPassword !== data.newPassword, {
        message: 'New password must be different than old password',
        path: ['newPassword'],
    })
    .refine((data) => data.newPassword === data.confirmNewPassword, {
        message: 'Passwords do not match!',
        path: ['confirmNewPassword'],
    });

export const ChangePasswordTab = () => {

    const auth = useAuth<AuthContextType>();
    const methods = useForm({
        resolver: zodResolver(schema),
        defaultValues: { ...auth.user }
    });

    return (<>
        <Card>
            <CardContent>

                <Form methods={methods} onSubmit={methods.handleSubmit(console.log)} >
                    <Stack gap={1.5}>
                        <PasswordFormField name="oldPassword" label="Old Password" required disabled />
                        <PasswordFormField name="newPassword" label="New Password" required disabled/>
                        <PasswordFormField name="confirmNewPassword" label="Confirm New Password" required disabled/>

                        <Stack direction="row" justifyContent="end">
                            <Button
                                type="submit"
                                variant="contained"
                                disabled
                            >
                                Save
                            </Button>
                        </Stack>
                    </Stack>
                </Form>

            </CardContent>
        </Card>
    </>)
}