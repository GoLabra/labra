
"use client";

import Bars3Icon from '@heroicons/react/24/outline/Bars3Icon';
import { Box, IconButton, SvgIcon, Theme, useMediaQuery } from '@mui/material';
import { useState } from 'react';
// import { Divide as HamburgerMenu } from 'hamburger-react';

interface HamburgerProps {
    onNavOpen: () => void;
}
export const Hamburger = (props: HamburgerProps) => {

    const { onNavOpen } = props;
    const [isOpen, setIsOpen] = useState(false);

    const mdDown = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

    return <>
        {mdDown && (
            <Box>
                <IconButton
                    color="inherit"
                    onClick={onNavOpen}
                >
                    <SvgIcon
                        color="action"
                        fontSize="small"
                    >
                        <Bars3Icon />
                    </SvgIcon>
                </IconButton>
            </Box>
        )}</>;
};
