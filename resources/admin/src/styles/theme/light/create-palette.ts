import type { ColorSystemOptions, PaletteOptions } from '@mui/material';
import { common } from '@mui/material/colors';
import { alpha, darken } from '@mui/material/styles';
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
      active: '#516377',
      disabled: alpha('#516377', 0.38),
      disabledBackground: alpha('#516377', 0.12),
      focus: alpha('#516377', 0.16),
      hover: alpha('#516377', 0.04),
      selected: alpha('#516377', 0.12)
    },
    background: {
      default: '#f5f7f9',
      paper: '#fdfdfd'
    },
    divider: '#b8bfcd',
    error,
    info,
    mode: 'light',
    neutral,
    primary: getPrimary(colorPreset),
    secondary: {
		light: '#6e7986',
		main: '#66707c',
		dark: '#565e68',
		contrastText: '#000000'
	},
    success,
    text: {
      primary: '#283543',
      secondary: '#66707c',
      disabled: alpha(neutral[900], 0.38)
    },
    warning,

    Switch: {

    }, 
    Avatar: {
        defaultBg: 'var(--mui-palette-text-primary)'
    }
  };
};
