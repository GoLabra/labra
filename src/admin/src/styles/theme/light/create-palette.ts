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
      active: '#516377',
      disabled: alpha('#516377', 0.38),
      disabledBackground: alpha('#516377', 0.12),
      focus: alpha('#516377', 0.16),
      hover: alpha('#516377', 0.04),
      selected: alpha('#516377', 0.12)
    },
    background: {
      default: '#f3f4f4',
      paper: '#ffffff'
    },
    divider: '#dce3e3',
    error,
    info,
    mode: 'light',
    neutral,
    primary: getPrimary(colorPreset),
    secondary: getPrimary(colorPreset),
    success,
    text: {
      primary: '#283543',
      secondary: '#57616c',
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
