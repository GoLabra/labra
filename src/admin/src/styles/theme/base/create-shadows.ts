import type { Shadows } from '@mui/material';
import { PaletteOptions } from '@mui/material';
import { alpha } from '@mui/material/styles';
import { neutral } from '../colors';

interface Config {
  palette: PaletteOptions;
}

export const createShadows = (): Shadows => {
  const layer1Color = alpha(neutral![200], 0.5);
  const layer2Color = alpha(neutral![800], 0.08);

  return [
    'none',
    `0px 0px 1px ${layer1Color}, 0px 1px 2px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 1px 3px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 3px 6px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 3px 8px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 4px 12px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 7px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 10px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 12px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 13px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 14px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 15px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 16px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 17px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 18px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 6px 20px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 28px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 28px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 30px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 34px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 32px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 34px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 36px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 38px ${layer2Color}`,
    `0px 0px 1px ${layer1Color}, 0px 10px 40px ${layer2Color}`
  ];
};
