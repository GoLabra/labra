import { getInitials, stringToColor } from "@/lib/utils/stringToColor";
import { Box, styled, Typography, Avatar as MuiAvatar } from "@mui/material";
import { useMemo } from "react";

export const AvatarShell = styled(Box)({
    padding: '8px',
    margin: 'auto',
    width: 124,
    height: 124,
    overflow: 'hidden',
    borderRadius: '100%',
    border: '2px dashed color-mix(in srgb, var(--mui-palette-background-paper), var(--mui-palette-common-onBackground) 20%)',

}) as typeof Box;


interface AvatarProps {
    name?: string;
}
export const Avatar = (props: AvatarProps) => {
    const name = useMemo(() => getInitials(props.name ?? ''), [props.name]);

    return (
        <AvatarShell>
            <Box width={1} height={1} sx={{ 
                borderRadius: '100%', 
                display: 'flex', 
                alignItems: 'center', 
                justifyContent: 'center', 
                bgcolor: 'var(--mui-palette-Avatar-defaultBg)',
                position: 'relative',
                color: 'var(--mui-palette-background-default)',
                '&::before': {
                    content: '""',
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    bottom: 0,
                    backgroundColor: stringToColor(name),
                }
                
        }}>
                <Typography variant="caption" fontSize={40} fontWeight={400}>{name}</Typography>
            </Box>
        </AvatarShell>
    )
}

export const InlineAvatar = (props: AvatarProps) => {

    const name = useMemo(() => getInitials(props.name ?? ''), [props.name]);

    return (<MuiAvatar sx={{
        width: 24,
        height: 24,
        fontSize: '0.875rem',
        position: 'relative',
        '&::before': {
            content: '""',
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: stringToColor(name),
        }
    }} >{name || ' '}</MuiAvatar>)
}