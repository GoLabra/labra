'use client'

import { createTheme } from '@mui/material/styles';
import { createShadows } from './base/create-shadows';
import { createOptions as createDarkOptions } from './dark/create-options';
import { createOptions as createLightOptions } from './light/create-options';
import { createComponents } from './base/create-components';
import { createTypography } from './base/create-typography';
import { COLOR } from '@/config/CONST';

declare module '@mui/material/styles' {

    interface Palette {
        neutral: NeutralColors;
    }

    interface PaletteOptions {
        neutral?: NeutralColors;
    }

    interface PaletteColor {
        alpha4?: string;
        alpha8?: string;
        alpha12?: string;
        alpha30?: string;
        alpha50?: string;
    }

    interface TypeBackground {
        paper: string;
        default: string;
    }
}

export interface NeutralColors {
    50: string;
    100: string;
    200: string;
    300: string;
    400: string;
    500: string;
    600: string;
    700: string;
    800: string;
    900: string;
}

export type ColorPreset = 'blue' | 'green' | 'indigo' | 'purple';

interface ThemeConfig {
    colorPreset?: ColorPreset,
}

let config: ThemeConfig = {
    colorPreset: COLOR as ColorPreset,
};


export const theme = createTheme({
    colorSchemes: {
        light: createLightOptions({
            colorPreset: config.colorPreset
        }),
        dark: createDarkOptions({
            colorPreset: config.colorPreset
        }),
    },
    cssVariables: {
        colorSchemeSelector: 'class',
    },
    typography: createTypography(),
    components: createComponents(),
    shadows: createShadows()
});
