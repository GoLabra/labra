import type { PaletteColor } from '@mui/material';
import { NeutralColors } from './index';
import { alpha } from '@mui/material/styles';

const withAlphas = (color: PaletteColor): PaletteColor => {
  return {
    ...color,
    alpha4: alpha(color.main, 0.04),
    alpha8: alpha(color.main, 0.08),
    alpha12: alpha(color.main, 0.12),
    alpha30: alpha(color.main, 0.30),
    alpha50: alpha(color.main, 0.50)
  };
};

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

export const blue = withAlphas({
  //light: 'color-mix(in srgb, var(--mui-palette-primary-main), var(--mui-palette-common-onBackground) 30%)',
  light: '#3273dd',
  main: '#3273dd',
  dark: '#3273dd',
  //dark: 'color-mix(in srgb, var(--mui-palette-primary-main), var(--mui-palette-common-background) 30%)',
  contrastText: '#FFFFFF'
});

export const green = withAlphas({
  light: '#6CE9A6',
  main: '#12B76A',
  dark: '#027A48',
  contrastText: '#FFFFFF'
});

export const indigo = withAlphas({
  light: '#EBEEFE',
  main: '#635dff',
  dark: '#4338CA',
  contrastText: '#FFFFFF'
});

export const purple = withAlphas({
  light: '#F4EBFF',
  main: '#9E77ED',
  dark: '#6941C6',
  contrastText: '#FFFFFF'
});

export const secondary = withAlphas({
  light: '#7f838e',
  main: '#7f838e',
  dark: '#7f838e',
  contrastText: '#FFFFFF'
});

export const success = withAlphas({
  light: '#4bbf73',
  main: '#4bbf73',
  dark: '#4bbf73',
  contrastText: '#FFFFFF'
});

export const info = withAlphas({
  light: '#1f9bcf',
  main: '#1f9bcf',
  dark: '#1f9bcf',
  contrastText: '#FFFFFF'
});

export const warning = withAlphas({
  light: '#cc8b37',
  main: '#cc8b37',
  dark: '#cc8b37',
  contrastText: '#FFFFFF'
});

export const error = withAlphas({
  light: '#bd464a',
  main: '#bd464a',
  dark: '#bd464a',
  contrastText: '#FFFFFF'
});
