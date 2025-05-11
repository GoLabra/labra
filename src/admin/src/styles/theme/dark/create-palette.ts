import type { ColorSystemOptions, PaletteOptions } from '@mui/material';
import { common } from '@mui/material/colors';
import { alpha } from '@mui/material/styles';
import { secondary, error, info, neutral, success, warning } from '../colors';
import type { ColorPreset } from '../index';
import { getPrimary } from '../utils';

interface Config {
  colorPreset?: ColorPreset;
}

export const createPalette = (config: Config): ColorSystemOptions['palette'] => {
  const { colorPreset } = config;

  return {
    action: {
      active: '#e9ecef',
      disabled: alpha('#e9ecef', 0.38),
      disabledBackground: alpha('#e9ecef', 0.12),
      focus: alpha('#e9ecef', 0.16),
      hover: alpha('#e9ecef', 0.04),
      selected: alpha('#e9ecef', 0.12)
    },
    background: {
      default: '#151b26',
      paper: '#282f3d'
    },
    divider: '#69717d36',
    error,
    info,
    mode: 'dark',
    neutral: neutral,
    primary: getPrimary(colorPreset),
    secondary,
    success,
    text: {
      primary: '#e9ecef',
      secondary: '#e9ecefA0',
      disabled: alpha(common.white, 0.38)
    },
    warning, 

    Switch: {
        
    },
    Avatar: {
        defaultBg: 'var(--mui-palette-text-primary)'
    }
  };
};
