import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import { Box, Stack } from "@mui/material";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {
    return (<Stack
        direction="row">
        <Box
            sx={{ width: '100%' }}>
            {props.children}
        </Box>
    </Stack>
    )
});

export default Layout