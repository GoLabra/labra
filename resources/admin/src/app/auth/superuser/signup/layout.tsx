
import ContentManagernEntities from "@/features/content-manager/content-manager-entities";
import { Box, Container, Stack } from "@mui/material";
import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";
import AppLayout from "@/shared/layouts/appLayout";
import AuthLayout from "@/shared/layouts/authLayout";

const Layout = (props: { children: React.ReactNode }) => {

    return (<AuthLayout>
        <Container maxWidth="xs">
            {props.children}
        </Container>
    </AuthLayout>)
}

export default Layout