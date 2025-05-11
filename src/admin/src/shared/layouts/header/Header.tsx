"use client";

import React, { ReactNode, useReducer, useState } from 'react';
import { AppBar, Box, IconButton, Stack, SvgIcon, Theme, Toolbar, useMediaQuery } from '@mui/material';
import { Hamburger } from './Hamburger/Hamburger';
import { ThemeSwitcher } from './ThemeSwitcher.tsx/ThemeSwitcher';
import { Notifications } from './Notification/Notification';
import { Topnav } from './Topnav/top-nav';
import { NavUser } from './User/user';
import { SkipLinks } from '@/shared/layouts/header/skip-links/src/skip-links';
import { Logo } from '@/shared/components/logo';
import { LogoDog } from '@/shared/components/logo.dog';

const TOP_NAV_HEIGHT: number = 64;


export interface HeaderProps {
    hasLeftPanel?: boolean;
    onNavOpen?: () => void;
    nav?: ReactNode,
    rightItems?: ReactNode[]
}
export const Header = (props: HeaderProps) => {
    const { onNavOpen } = props;

    const mdDown = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));


    return (
        <>
            <AppBar position="sticky" sx={{
                  gridArea: 'header' 
                }}>

                <Toolbar disableGutters variant='dense'
                    sx={{
                        height: TOP_NAV_HEIGHT
                    }}>
                    <Stack direction="row" sx={{ width: {md: '240px', lg: '240px', xl: '240px'} }}>
                        <Box px={3}>
                            <Logo />
                        </Box>

                        <SkipLinks />

                        {(onNavOpen && (props.hasLeftPanel || props.nav)) && <Hamburger onNavOpen={onNavOpen} />}
                    </Stack>

                    {mdDown == false && (props.nav)}

                    <Box sx={{ flexGrow: 1 }} >
                    </Box>

                    <Stack
                        direction="row"
                        gap={2}
                        pr={3}>
                            {props.rightItems}

                        {/* <Notifications total={20} />
                        <ThemeSwitcher />
                        <NavUser /> */}

                        {/* <Messages total={15} />
						<UserAccount onClick={handleProfileMenuOpen} /> */}
                    </Stack>
                    
                </Toolbar>
            </AppBar>
        </>
    );
};

export const loggedInHeaderProps = {
    nav: <Topnav key="top-nav" /> as ReactNode,
    rightItems: [
        <Notifications key="notifications" />,
        <NavUser key="user" />
    ] as ReactNode[]
}