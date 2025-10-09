
import ContentManagernEntities from "@/features/content-manager/content-manager-entities";
import { Box, Stack } from "@mui/material";
import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import AppLayout from "@/shared/layouts/appLayout";
import { CurrentEntityProvider } from "@/hooks/use-current-entity";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {

    return (
        <CurrentEntityProvider>
            <AppLayout sideChildren={<ContentManagernEntities />}>
                <Stack
                    direction="row">
                    <Box
                        sx={{ width: '100%' }}>
                        {props.children}
                    </Box>
                </Stack>
            </AppLayout>
        </CurrentEntityProvider>
    )
})

export default Layout