import { withAuthGuard } from "@/core-features/auth/hocs/with-auth-guard";

const Layout = withAuthGuard((props: { children: React.ReactNode }) => {
    return props.children
});

export default Layout