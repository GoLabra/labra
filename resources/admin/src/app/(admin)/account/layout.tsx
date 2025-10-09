import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import AppLayout from "@/shared/layouts/appLayout";
import { Box, Stack } from "@mui/material";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {

    return (
        <AppLayout>
            <Stack direction="row" width={1} >
                <Box
                    sx={{ width: '100%' }}>
                    {props.children}
                </Box>
            </Stack>
        </AppLayout>
    )
})

export default Layout