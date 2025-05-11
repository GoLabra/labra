import type { ColorSystemOptions, ThemeOptions } from '@mui/material';
import type { ColorPreset } from '../index';
// import { createComponents } from './create-components';
import { createPalette } from './create-palette';
// import { createShadows } from './create-shadows';

interface Config {
    colorPreset?: ColorPreset;
}

export const createOptions = (config: Config): ColorSystemOptions => {
    const { colorPreset } = config;
    const palette = createPalette({ colorPreset });
    //   const components = createComponents({ palette });
    //   const shadows = createShadows();

    return {
        palette,
        //shadows
    };
};
