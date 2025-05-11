import ChevronDownIcon from "@heroicons/react/24/outline/ChevronDownIcon";
import ChevronRightIcon from "@heroicons/react/24/outline/ChevronRightIcon";
import { Box, ButtonBase, Collapse, SvgIcon } from "@mui/material";
import { ReactNode, useCallback, useState } from "react";
import NextLink from 'next/link';

interface SideNavItemProps {
    active?: boolean;
    children?: ReactNode;
    collapse?: boolean;
    depth?: number;
    externalLink?: boolean;
    icon?: ReactNode;
    openImmediately?: boolean;
    path?: string;
    pinned?: boolean;
    title: string;
}

export const SideNavItem: React.FC<SideNavItemProps> = (props: SideNavItemProps) => {
    const {
        active = false,
        children,
        collapse = false,
        depth = 0,
        externalLink = false,
        icon,
        openImmediately = true,
        path,
        title
    } = props;
    const [open, setOpen] = useState<boolean>(openImmediately);

    const handleToggle = useCallback(
        (): void => {
            setOpen((prevOpen) => !prevOpen);
        },
        []
    );

    // Branch
    if (children) {
        return (
            <li>
                <ButtonBase
                    onClick={handleToggle}
                    sx={{
                        alignItems: 'center',
                        borderRadius: 1,
                        color: 'text.secondary',
                        display: 'flex',
                        fontFamily: (theme) => theme.typography.fontFamily,
                        fontSize: 12,
                        fontWeight: 400,
                        justifyContent: 'flex-start',
                        py: '8px',
                        pr: '8px',
                        mt: '20px',
                        mb: '10px',
                        textAlign: 'left',
                        whiteSpace: 'nowrap',
                        width: '100%'
                    }}
                >

                    <Box
                        component="span"
                        sx={{
                            // color: depth === 0 ? 'text.primary' : 'text.secondary',
                            flexGrow: 1,
                            fontSize: 12,
                            mx: '10px',
                            transition: 'opacity 250ms ease-in-out',
                            // ...(active && {
                            //   color: 'primary.main'
                            // }),
                            ...(collapse && {
                                opacity: 0
                            })
                        }}
                    >
                        {title.toUpperCase()}
                    </Box>
                    <SvgIcon
                        sx={{
                            color: 'neutral.500',
                            fontSize: 16,
                            transition: 'opacity 250ms ease-in-out',
                            ...(collapse && {
                                opacity: 0
                            })
                        }}
                    >
                        {open ? <ChevronDownIcon /> : <ChevronRightIcon />}
                    </SvgIcon>
                </ButtonBase>
                <Collapse
                    in={!collapse && open}
                    unmountOnExit
                >
                    {children}
                </Collapse>
            </li>
        );
    }

    // Leaf

    const linkProps = path
        ? externalLink
            ? {
                component: 'a',
                href: path,
                target: '_blank'
            }
            : {
                component: NextLink,
                href: path
            }
        : {};

    return (
        <li>
            <ButtonBase
                sx={{
                    alignItems: 'center',
                    color: 'text.secondary',
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
                        backgroundColor: 'background.default'
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
                            //color: 'primary.main'
                        }),
                        ...(collapse && {
                            opacity: 0
                        })
                    }}
                >
                    {title}
                </Box>

                {externalLink && (
                    <SvgIcon
                        sx={{
                            color: 'neutral.500',
                            fontSize: 18,
                            transition: 'opacity 250ms ease-in-out',
                            ...(collapse && {
                                opacity: 0
                            })
                        }}
                    >
                        {/* <ArrowTopRightOnSquareIcon /> */}
                    </SvgIcon>
                )}
            </ButtonBase>
        </li>
    );
};
