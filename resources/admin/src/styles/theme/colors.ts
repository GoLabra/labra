import type { PaletteColor } from '@mui/material';
import { NeutralColors } from './index';
import { alpha } from '@mui/material/styles';

export const neutral: NeutralColors = {
  50: '#CBD3E2',
  100: '#B8BFCD',
  200: '#A5ACB9',
  300: '#9298A4',
  400: '#7F8590',
  500: '#6C717B',
  600: '#595E67',
  700: '#464A52',
  800: '#33373E',
  900: '#202329'
};

export const blue = {
  //light: 'color-mix(in srgb, var(--mui-palette-primary-main), var(--mui-palette-common-onBackground) 30%)',
  light: '#4475b2',
  main: '#4475b2',
  dark: '#4475b2',
  //dark: 'color-mix(in srgb, var(--mui-palette-primary-main), var(--mui-palette-common-background) 30%)',
  contrastText: '#FFFFFF'
};

export const green = {
  light: '#6CE9A6',
  main: '#12B76A',
  dark: '#027A48',
  contrastText: '#FFFFFF'
};

export const indigo = {
  light: '#EBEEFE',
  main: '#635dff',
  dark: '#4338CA',
  contrastText: '#FFFFFF'
};

export const purple = {
  light: '#F4EBFF',
  main: '#9E77ED',
  dark: '#6941C6',
  contrastText: '#FFFFFF'
};

export const success = {
  light: '#4bbf73',
  main: '#4bbf73',
  dark: '#4bbf73',
  contrastText: '#FFFFFF'
};

export const info = {
  light: '#1f9bcf',
  main: '#1f9bcf',
  dark: '#1f9bcf',
  contrastText: '#FFFFFF'
};

export const warning = {
  light: '#cc8b37',
  main: '#cc8b37',
  dark: '#cc8b37',
  contrastText: '#FFFFFF'
};

export const error = {
  light: '#bd464a',
  main: '#bd464a',
  dark: '#bd464a',
  contrastText: '#FFFFFF'
};
