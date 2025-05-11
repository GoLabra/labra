import type { ThemeOptions } from '@mui/material';
import type { ColorPreset } from '../index';
// import { createComponents } from './create-components';
import { createPalette } from './create-palette';
// import { createShadows } from './create-shadows';

interface Config {
    colorPreset?: ColorPreset;
}

export const createOptions = (config: Config): ThemeOptions => {
    const { colorPreset } = config;
    const palette = createPalette({ colorPreset });
    //   const components = createComponents({ palette });
    //   const shadows = createShadows({ palette });

    return {
        palette,
        components: {
            
        }
        //shadows
    };
};
