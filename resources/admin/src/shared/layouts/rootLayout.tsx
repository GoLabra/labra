"use client";

import { Header, HeaderProps } from "@/shared/layouts/header/Header";
import MobileSideNav from "@/shared/layouts/side-nav/mobile-side-nav";
import { Box, BoxProps, Divider, Stack, Theme, styled, useMediaQuery } from "@mui/material";
import { PropsWithChildren, ReactNode, useMemo, useState } from "react";
import SideNav from "./side-nav/side-nav";

export const Container = styled(Box)(({ theme }) => ({
    display: 'grid',

    gridTemplateAreas:`
      "header header"
      "nav content"`,
  
    gridTemplateColumns: `auto 1fr`,
    gridTemplateRows: 'auto 1fr',
    gridHap: '10px',
    minHeight: '100vh',

    [theme.breakpoints.down('md')]: {
        gridTemplateColumns: '1fr',
        gridTemplateAreas:`
            "header "
            "content"`,
            }
}));


interface LeftSidePanelProps {
    mobileNavOpen: boolean;
    setMobileNavOpen: (open: boolean) => void;
    mobileNav?: ReactNode;
}
export const LeftSidePanel = (props: PropsWithChildren<LeftSidePanelProps>) => {
    const mdDown = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

    if (mdDown) {
        return <MobileSideNav
            open={props.mobileNavOpen}
            onClose={() => props.setMobileNavOpen(false)}>
            <Stack>
                {props.mobileNav}
                <Divider />
                {props.children}
            </Stack>
        </MobileSideNav>
    }

    if (props.children) {
        return <SideNav>{props.children}</SideNav>
    }

    return null;
}

export const ContentLayout = styled(Box)({
    width: '100%',
    gridArea: 'content'
});