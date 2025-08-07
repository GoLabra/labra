import type { ColorSystemOptions, PaletteOptions } from '@mui/material';
import { common } from '@mui/material/colors';
import { alpha } from '@mui/material/styles';
import { error, info, neutral, success, warning } from '../colors';
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
      default: '#14171e',
      paper: '#1b212c'
    },
    divider: '#3c434c',
    error,
    info,
    mode: 'dark',
    neutral: neutral,
    primary: getPrimary(colorPreset),
    secondary: {
		light: '#a8b4c4',
		main: '#9198a1',
		dark: '#828b97',
		contrastText: '#FFFFFF'
	},
    success,
    text: {
      primary: '#e9ecef',
      secondary: '#9198a1',
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
