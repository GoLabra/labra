import ChevronDownIcon from "@heroicons/react/24/outline/ChevronDownIcon";
import ChevronRightIcon from "@heroicons/react/24/outline/ChevronRightIcon";
import { Box, ButtonBase, Collapse, SvgIcon } from "@mui/material";
import { ReactNode, useCallback, useState } from "react";
import NextLink from 'next/link';
import { alpha, darken } from '@mui/material/styles';

interface SideNavLiteItemProps {
    active?: boolean;
    icon?: ReactNode;
    path?: string;
    title: string;
}

export const SideNavLiteItem: React.FC<SideNavLiteItemProps> = (props: SideNavLiteItemProps) => {
    const {
        active = false,
        icon,
        path,
        title
    } = props;


    const linkProps = path
        ? {
            component: NextLink,
            href: path
        }
        : {};

    return (
        <li>
            <ButtonBase
                sx={{
                    alignItems: 'center',
                    color: 'text.primary',
                    borderRadius: 1,
                    display: 'flex',
                    fontFamily: (theme) => theme.typography.fontFamily,
                    fontSize: 14,
                    fontWeight: 400,
                    justifyContent: 'flex-start',
                    p: '6px 8px',
                    my: '2px',
                    textAlign: 'left',
                    whiteSpace: 'nowrap',
                    width: '100%',

                    ...(active && {
                        //backgroundColor: 'background.paper',
                        //backgroundColor: (theme) => darken(theme.palette.primary.alpha8!, 1),
                        color: 'primary.main'
                    }),
                }}
                {...linkProps}
            >
                <Box
                    component="span"
                    sx={{
                        alignItems: 'center',

                        display: 'inline-flex',
                        flexGrow: 0,
                        flexShrink: 0,
                        height: 18,
                        justifyContent: 'center',
                        width: 18,
                        '>svg ': {
                            fontSize: '18px'
                        }
                    }}
                >
                    {icon}
                </Box>

                <Box
                    component="span"
                    sx={{

                        flexGrow: 1,
                        mx: '6px',
                        ml: '12px',
                        transition: 'opacity 250ms ease-in-out',
                        ...(active && {
                            color: 'primary.main'
                        }),

                    }}
                >
                    {title}
                </Box>

            </ButtonBase>
        </li>
    );
};
