'use client'

import { Box, List, ListItem, ListItemIcon, ListItemText, ListItemButton, styled, IconButton, SvgIcon, ButtonProps, ListItemButtonProps, alpha, ListItemTextProps, Skeleton, Stack } from "@mui/material"
import NextLink from 'next/link';
import { PropsWithChildren, useMemo } from "react";

export const TreeListItem = styled(ListItem)(({ theme }) => {
    return [
        {
            paddingLeft: '12px',
            borderLeftStyle: 'solid',
            borderLeftWidth: 2,
            '&:before': {
                content: '""',
                display: 'block',
                position: 'absolute',
                top: '-14px',
                left: '-2px',
                width: '12px',
                height: '35px',
                borderStyle: 'solid',
                borderWidth: '0 0 2px 2px',
                borderBottomLeftRadius: '8px'
            },
            '&:last-child': {
                borderLeftColor: 'transparent'
            },

            '&.branch-auto-top:before': {
                top: 'inherit',
                transform: 'translateY(-50%)'
            }
        },
        theme.applyStyles('light', {
            borderLeftColor: 'var(--mui-palette-neutral-100)',
            '&:before': {
                borderColor: 'var(--mui-palette-neutral-100)',
            }
        }),
        theme.applyStyles('dark', {
            borderLeftColor: 'var(--mui-palette-background-paper)',
            '&:before': {
                borderColor: 'var(--mui-palette-background-paper)',
            }
        }),
    ]
});

interface TreeListProps {
    children: React.ReactNode
}
export const TreeList = (props: TreeListProps) => {

    const { children } = props;

    return (<List sx={{ padding: '5px 0 0 16px', overflow: 'hidden', position: "relative" }}>
        {children}
    </List>)
}

interface TreeListItemTextProps extends ListItemTextProps {
    label: string;
    icon?: React.ReactNode;
}
export const TreeListItemText = (props: PropsWithChildren<TreeListItemTextProps>) => {
    const { label, icon, children, ...other } = props;

    return (
        <TreeListItem
            disablePadding >
            <Stack width={1}>

                <ListItemText primary={label} {...other} slotProps={{
                    primary: {
                        sx: {
                            fontSize: '0.9rem',
                            paddingLeft: '10px'
                        },
                        ...other.slotProps?.primary
                    }
                }} />

                {icon}
            </Stack>

            {children}

        </TreeListItem>)
}


interface TreeListItemButtonProps extends ListItemButtonProps {
    label: string;
    icon?: React.ReactNode;
}
export const TreeListItemButton = (props: PropsWithChildren<TreeListItemButtonProps & ListItemButtonProps>) => {
    const { label, icon, children, ...other } = props;

    const ListItemButtonStyled = styled(ListItemButton)<ListItemButtonProps>(({ theme }) => ({
        padding: '0 10px',
        borderRadius: 1
    }));

    return (
        <TreeListItem
            disablePadding >

            <ListItemButtonStyled
                dense
                {...other} >

                <ListItemText primary={label}
                    slotProps={{
                        primary: {
                            fontSize: '0.9rem',
                        }
                    }}
                />

                {icon}

            </ListItemButtonStyled>

            {children}


        </TreeListItem>)
}


interface TreeListItemNavigationProps extends ListItemButtonProps {
    label: string;
    path: string;
    icon?: React.ReactNode;
    active: boolean;
    externalLink?: boolean;
}
export const TreeListItemNavigation = (props: TreeListItemNavigationProps) => {
    const { label, path, icon, active, children, externalLink, ...other } = props;

    const ListItemButtonStyled = styled(ListItemButton)<ListItemButtonProps & { href: string }>(({ theme }) => ({
        marginTop: '5px',
        marginBottom: '5px',
        padding: '0px 10px',
        borderRadius: 1,
        ...(active && {
            // color: 'primary.main',
            backgroundColor: alpha(theme.palette.primary.main, 0.1),
        }),
    }));

    const linkProps = externalLink
        ? {
            component: 'a',
            href: path,
            target: '_blank'
        }
        : {
            component: NextLink,
            href: path
        } as any;

    return (
        <TreeListItem
            disablePadding >
            <ListItemButtonStyled
                {...other}
                {...linkProps}
                LinkComponent={NextLink}>
                <ListItemText
                    slotProps={{
                        primary: {
                            typography: {
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                fontSize: '0.9rem',
                                ...(active && {
                                    color: 'var(--mui-palette-primary-main)'
                                }),
                            }
                        }
                    }}
                    primary={label} />

                {icon}
            </ListItemButtonStyled>

            {children}

        </TreeListItem>)
}



export const TreeListItemSkeleton = () => {
    const textSx = useMemo(() => ({
        paddingLeft: '5px'
    }), [])

    return (
        <TreeListItem>
            <ListItemText primary={<Skeleton variant="text" />} sx={textSx}>
            </ListItemText>
        </TreeListItem>)
}
