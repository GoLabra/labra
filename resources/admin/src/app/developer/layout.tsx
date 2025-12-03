
import { Box, Stack } from "@mui/material";
import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import DeveloperLeftSide from "@/features/developer/developer-left-side";
import DevLayout from "@/shared/layouts/devLayout";
import { CurrentEntityProvider } from "@/hooks/use-current-entity";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {

    return (
        <CurrentEntityProvider>
            <DevLayout sideChildren={<DeveloperLeftSide />}>
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