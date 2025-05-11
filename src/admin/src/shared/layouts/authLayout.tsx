"use client";

import { Container, ContentLayout } from "./rootLayout";
import { ThemeSwitcher } from "./header/ThemeSwitcher.tsx/ThemeSwitcher";
import { Box, Stack } from "@mui/material";
import { Header } from "./header/Header";

interface Props {
    children: React.ReactNode;
    sideChildren?: React.ReactNode;
}
export default function AuthLayout(props: Props) {

    return (<Container
        sx={{
            gridTemplateAreas: `
                                "header header"
                                "art content"`,
            gridTemplateColumns: `1fr 1fr`,
            gridTemplateRows: 'auto 1fr',
        }}>
        <Header
            rightItems={[<ThemeSwitcher key="theme-switcher" />]}
            hasLeftPanel={!!props.sideChildren}
        />
        <Box sx={{
            gridArea: 'art',
            background: 'url(/graph-background.png)',
            backgroundRepeat: 'no-repeat',
            backgroundSize: '100%',
            backgroundPositionY: 'bottom',
            borderRight: '1px solid var(--mui-palette-divider)',
        }} />

        <ContentLayout id="main">
            <Stack direction="row" width={1} height={1} alignItems="center">
                {props.children}
            </Stack>
        </ContentLayout>

    </Container>)
}
