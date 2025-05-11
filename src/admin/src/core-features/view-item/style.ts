import { styled } from "@mui/material";

export const ArrowRoot = styled('div')(({

    position: 'absolute',
    top: '-9px',
    left: '50%',
    transform: 'translateX(-50%)',
    width: 0,
    height: 0,
    borderLeft: '10px solid transparent',
    borderRight: '10px solid transparent',
    borderBottom: '10px solid var(--mui-palette-background-paper)',
}
));
