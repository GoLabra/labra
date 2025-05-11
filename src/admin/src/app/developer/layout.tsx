
import { Box, Stack } from "@mui/material";
import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import DeveloperEntities from "@/features/developer/developer-entities";
import DevLayout from "@/shared/layouts/devLayout";
import { CurrentEntityProvider } from "@/hooks/use-current-entity";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {

    return (
        <CurrentEntityProvider>
            <DevLayout sideChildren={<DeveloperEntities />}>
                <Stack
                    direction="row">
                    <Box
                        sx={{ width: '100%' }}>
                        {props.children}
                    </Box>
                </Stack>
            </DevLayout>
        </CurrentEntityProvider>
    )
})

export default Layout