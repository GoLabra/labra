import type { FC } from 'react';
import PropTypes from 'prop-types';
import { useTheme } from '@mui/material/styles';
import { Box } from '@mui/material';


interface LogoProps {
    color?: any;
}
export const Logo: FC<LogoProps> = (props) => {
    const { color: colorProp } = props;
    const theme = useTheme();

    const color = 'var(--mui-palette-text-secondary)';

    return (
        <Box tabIndex={0} sx={{
            outlineWidth: '0px',
            height: '40px'
        }}>
            <svg width="40" height="40" viewBox="0 0 123 116" fill="none" xmlns="http://www.w3.org/2000/svg">



        <path d="M50 49.5V116L70 105L50 49.5Z" fill={color} />
<path d="M66.5 90L50 44H111L66.5 90Z" fill={color} />
<path d="M4 1L48 41V0L4 1Z" fill={color} />
<path d="M50 2V42H88L50 2Z" fill={color} />
<path d="M8 7L0 9L25.5 23L8 7Z" fill={color} opacity="0.5"/>
<path d="M114 13V23H125L114 13Z" fill={color} />
<path d="M103.5 11L81.5 32.5L90.5 42H112V11H103.5Z" fill={color} />


            </svg>
        </Box>
    );
};