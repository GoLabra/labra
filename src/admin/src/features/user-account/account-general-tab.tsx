import { AuthContextType } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { AvatarUploadField } from "@/core-features/dynamic-form/form-fields/AvatarUploadField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { UploadAvatar } from "@/shared/components/upload/upload-avatar";
import { zodResolver } from "@hookform/resolvers/zod";
import { Avatar, Box, Button, Card, CardContent, Grid, Stack, Typography } from "@mui/material"
import { useForm } from "react-hook-form";
import { z } from "zod";


const avatarSchema = z.object({
    avatar: z.custom().transform((data, ctx) => {
        const hasFile = data instanceof File || (typeof data === 'string' && !!data.length);
  
        if (!hasFile) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            message: 'Avatar is required!',
          });
          return null;
        }
  
        return data;
      }),
});

const schema = z.object({
    name: z.string().min(1, 'Name is required'),
    email: z.string().min(1, 'Name is required').email(), 
    firstName: z.string().min(1, 'First Name is required'),
    lastName: z.string().min(1, 'Last Name is required'),
});

export const AccountGeneralTab = () => {

    
    const auth = useAuth<AuthContextType>();
    const methods = useForm({
        resolver: zodResolver(schema),
        defaultValues: {...auth.user}
    });

    const avatarMethods = useForm({
        resolver: zodResolver(avatarSchema),
    });
    
    return (<>

        <Box>

            <Grid container spacing={2} className="andrei">
                <Grid item xs={12} sm={4}>
                    <Card>
                        <CardContent>
                            <Stack>
                                <br /><br /><br />
                                <Form methods={avatarMethods} onSubmit={methods.handleSubmit(console.log)} >
                                    <AvatarUploadField name="avatar" disabled />
                                </Form>

                                <br />
                                {/* <Typography color="text.secondary" sx={{ fontSize: '0.8rem' }}>
                                    Allowed *.jpeg, *.jpg, *.png, *.gif
                                    max size of 3 Mb
                                </Typography> */}
                                <br /><br /><br />
                                <br /><br /><br />
                            </Stack>

                        </CardContent>
                    </Card>
                </Grid>
                <Grid item xs={12} sm={8}>
                    <Card>
                        <CardContent>

                            <Form methods={methods} onSubmit={methods.handleSubmit(console.log)} >
                                <Stack gap={1.5}>
                                    <TextShortFormField name="name" label="Name" required disabled/>
                                    <TextShortFormField name="email" label="Email" required disabled/>
                                    <TextShortFormField name="firstName" label="First Name" required disabled/>
                                    <TextShortFormField name="lastName" label="Last Name" required disabled/>

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
                </Grid>
            </Grid>

        </Box>
    </>)
}