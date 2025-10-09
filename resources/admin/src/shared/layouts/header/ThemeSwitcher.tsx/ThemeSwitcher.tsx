"use client";

import { IconButton, SvgIcon, useColorScheme } from '@mui/material';
import MoonIcon from '@heroicons/react/24/outline/MoonIcon';
import SunIcon from '@heroicons/react/24/outline/SunIcon';

export const ThemeSwitcher = () => {

    const { mode, setMode } = useColorScheme();
    const currentMode = mode || 'dark';

    return (
        <IconButton
            color="inherit"
            onClick={() => {
                setMode(mode === 'dark' ? 'light' : 'dark')
            }}
            aria-label="Toggle theme"
        >
            <SvgIcon
                color="action"
                fontSize="small"
            >
                {currentMode === 'dark' ? <SunIcon /> : <MoonIcon />}
            </SvgIcon>
        </IconButton>
    );
};
