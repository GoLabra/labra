import { usePopover } from "@/hooks/use-popover";
import ArrowRightOnRectangleIcon from "@heroicons/react/24/outline/ArrowRightOnRectangleIcon";
import UserIcon from "@heroicons/react/24/outline/UserIcon";
import WrenchScrewdriverIcon from "@heroicons/react/24/outline/WrenchScrewdriverIcon";
import { Button, ButtonProps, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Menu, Popover, SvgIcon, Switch } from "@mui/material";
import WifiIcon from '@mui/icons-material/Wifi';
import React, { useMemo } from "react";
import { ResponsiveButton } from "@/styles/button.responsive";


interface MenuButtonProps {
    text?: React.ReactNode;
    disabled?: boolean;
    children: React.ReactNode;
    slots?: {
        button?: React.ReactNode;
    },
    slotProps?: {
        buttonProps?: ButtonProps
    }
}
export const MenuButton = (props: MenuButtonProps) => {

    const { disabled, text, children, slotProps } = props;
    const popover = usePopover<HTMLButtonElement>();

    const buttonProps = useMemo(() => ({
        onClick: popover.handleOpen,
        ref: popover.anchorRef,
        disabled: disabled,
        'aria-controls': popover.open ? 'basic-menu' : undefined,
        'aria-haspopup': "true" as "true",
        'aria-expanded': popover.open ? 'true' : undefined as "true" | undefined
    }), []);

    return (
        <>

            {props?.slots?.button 
            ? React.cloneElement(props.slots.button as React.ReactElement, {
                ...buttonProps
            }) 
            : (<ResponsiveButton
                variant="text"
                size="medium"
                {...buttonProps}
                {...slotProps?.buttonProps}
            >
                <span className="button-text">{text}</span></ResponsiveButton>)}

            <Menu
                id="basic-menu"
                anchorEl={popover.anchorRef.current}
                anchorOrigin={{
                    horizontal: 'right',
                    vertical: 'bottom'
                }}
                transformOrigin={{
                    horizontal: 'right',
                    vertical: 'top'
                  }}
                open={popover.open}
                onClose={popover.handleClose}
                MenuListProps={{
                    'aria-labelledby': 'basic-button',
                }}>

                {React.Children.map(children, (child) => {
                    if (React.isValidElement(child) && child.props["data-autoclose"]) {
                        return React.cloneElement(child as React.ReactElement, {
                            onClick: (event: React.MouseEvent<HTMLLIElement>) => {
                                child.props.onClick?.(event); // Preserve original onClick
                                popover.handleClose(); // Auto-close if 'data-autoclose' is set
                            },
                        });
                    }
                    return child;
                })}

            </Menu>
        </>
    )
}