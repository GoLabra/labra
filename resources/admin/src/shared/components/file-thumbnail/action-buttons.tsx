import { useTheme } from '@mui/material/styles';
import ButtonBase from '@mui/material/ButtonBase';
import IconButton from '@mui/material/IconButton';
import CloseIcon from '@mui/icons-material/Close';
import HighlightOffIcon from '@mui/icons-material/HighlightOff';

//import { bgBlur, varAlpha } from 'src/theme/styles';

// import { Iconify } from '../iconify';

// ----------------------------------------------------------------------
interface DownloadButtonProps {
    sx?: any;
    onClick?: () => void;
    className?: string;
}
export function DownloadButton({ sx, ...other }: DownloadButtonProps) {
    const theme = useTheme();

    return (
        <ButtonBase
            sx={{
                p: 0,
                top: 0,
                right: 0,
                width: 1,
                height: 1,
                zIndex: 9,
                opacity: 0,
                position: 'absolute',
                color: 'common.white',
                borderRadius: 'inherit',
                transition: theme.transitions.create(['opacity']),
                '&:hover': {
                    //...bgBlur({ color: varAlpha(theme.vars.palette.grey['900Channel'], 0.64) }),
                    opacity: 1,
                },
                ...sx,
            }}
            {...other}
        >
            <HighlightOffIcon fontSize="inherit" />
        </ButtonBase>
    );
}

// ----------------------------------------------------------------------

interface RemoveButtonProps {
    sx?: any;
    className?: string;
    onClick?: () => void;
}
export function RemoveButton({ sx, ...other }: RemoveButtonProps) {
    return (
        <IconButton
            size="small"
			
            sx={{
                p: 0.35,
                top: 4,
                right: 4,
                position: 'absolute',
                color: 'common.white',
                //bgcolor: (theme) => varAlpha(theme.vars.palette.grey['900Channel'], 0.48),
                //'&:hover': { bgcolor: (theme) => varAlpha(theme.vars.palette.grey['900Channel'], 0.72) },
                ...sx,
            }}
            {...other}
        >
			<HighlightOffIcon fontSize="inherit" />
        </IconButton>
    );
}
