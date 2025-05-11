'use client'

import { AuthContextType, STORAGE_KEY } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { paths } from "@/lib/paths";
import { Copywrite } from "@/shared/components/copywrite";
import { Alert, Box, Button, Stack, TextField, Typography } from "@mui/material";
import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useMemo } from "react";
import { gql, useLazyQuery, useQuery } from "@apollo/client";
import { getJwtRole, getJwtSub } from "@/lib/utils/jwt";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { PasswordFormField } from "@/core-features/dynamic-form/form-fields/PasswordField";
import { SelectFormField } from "@/core-features/dynamic-form/form-fields/SelectField";
import { z } from "zod";
import { Avatar } from "@/shared/components/avatar";


const getMeRoles = gql`query getMeRoles($where:UserWhereInput!) {
    users(where:$where) {
          roles {
              id
              name
          }
      }
  }`


const schema = z.object({
    role: z.string().nonempty('Role is required'),
});

export default function AuthPage() {

    const auth = useAuth<AuthContextType>();
    const accessToken = globalThis.localStorage.getItem(STORAGE_KEY);
    const roles = useQuery<any>(getMeRoles, { variables: { where: { email: getJwtSub(accessToken) } },
                                            fetchPolicy: 'network-only' });

    const router = useRouter();
    const searchParams = useSearchParams();
    const returnTo = searchParams.get('returnTo') || undefined;
    

    const currentRole = useMemo(() => getJwtRole(globalThis.localStorage.getItem("accessToken")), []);

    const roleOptions = useMemo(() => roles.data?.users[0].roles?.map((i:any) => ({
        label: i.name,
        value: i.name
    })) ?? [], [roles.data?.users[0].roles]);

    const methods = useForm({
        resolver: zodResolver(schema),
        defaultValues: useMemo(() => ({ role: currentRole}), [currentRole])
    });

    const onSubmit = useCallback(async (data: any) => {
        await auth.changeRole(data.role);
        router.push(returnTo || paths.contentManager.index);
    }, [auth.signUp]);

    return (<>

        <Box sx={{ pt: 7, pb: 14 }}>

            {/* <Stack direction="row" alignItems="center" justifyContent="center" sx={{ ...avatarBorderStyle, mb: 7 }}>
                <Avatar sx={{ width: 65, height: 65, fontSize: '2rem', ...avatarProps.sx }} >{avatarProps.children}</Avatar>
            </Stack> */}

            <Avatar name={auth.user?.name ?? ''} />

            <Stack
                alignItems="center"
                direction="row"
                justifyContent="end"
                spacing={1}
                sx={{ mb: 3 }}
            >
                <Button
                    onClick={() => router.back()}
                >
                    Go Back
                </Button>
            </Stack>


            <Form methods={methods} onSubmit={methods.handleSubmit(onSubmit)} >
                <Stack gap={1.5}>
                    <SelectFormField name="role" label="Role" options={roleOptions} required />
                </Stack>

                <Button
                    type="submit"
                    fullWidth
                    sx={{ mt: 3 }}
                    variant="contained"
					disabled={methods.formState.isSubmitting}
                >
                    Login with role
                </Button>
            </Form>

            <Box textAlign="center" sx={{ mt: 3 }}>
                <Copywrite />
            </Box>

        </Box >
    </>
    )
}
