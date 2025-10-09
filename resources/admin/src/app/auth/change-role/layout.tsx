
import { Container } from "@mui/material";
import AuthLayout from "@/shared/layouts/authLayout";
import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {

    return (<AuthLayout>
        <Container maxWidth="xs">
            {props.children}
        </Container>
    </AuthLayout>)
});

export default Layout