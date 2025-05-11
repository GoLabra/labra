"use client";

import PropTypes from 'prop-types';
import { Button, Divider, ListItemIcon, Menu, MenuItem } from '@mui/material';
import { usePopover } from '../../hooks/use-popover';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';

const isCriticStyle = {
    color: 'error.main',
    '&:hover': {
        backgroundColor: 'error.lighter'
    }
}

interface ActionBase {
    type: 'action' | 'separator';
}
type Action = (ActionBase & {
    type: 'action';
    handler?: () => void;
    disabled?: boolean;
    isCritic?: boolean
    label: string;
    icon?: JSX.Element;
}) | (ActionBase & {
    type: 'separator';
})

interface ActionsButtonProps {
    actions?: Action[];
    label?: string;
    variant?: 'text' | 'outlined' | 'contained';
    disabled?: boolean;
    color?: 'inherit' | 'primary' | 'secondary' | 'success' | 'error' | 'info' | 'warning'
}

export const ActionsButton = (props: ActionsButtonProps) => {
    const { actions = [], label, variant, disabled, color } = props;
    const popover = usePopover<HTMLButtonElement>();


    return (
        <>
            <Button
                size="medium"
                onClick={popover.handleOpen}
                ref={popover.anchorRef}

                variant={variant}
                disabled={disabled}
                color={color}

                aria-label="Actions dropdown"
                aria-haspopup="menu"
                aria-expanded={popover.open ? 'true' : undefined}
            >
                {label}
                <ArrowDropDownIcon />
            </Button>


            <Menu
                anchorEl={popover.anchorRef.current}
                anchorOrigin={{
                    horizontal: 'right',
                    vertical: 'bottom'
                }}
                open={popover.open}
                onClose={popover.handleClose}
                transformOrigin={{
                    horizontal: 'right',
                    vertical: 'top'
                }}
            >
                {actions.map((item, index) => {

                    if (item.type == 'separator') {
                        return <Divider key={`divider-${index}`} />
                    }

                    return (<MenuItem
                        key={item.label}
                        disabled={item.disabled === true}
                        sx={{
                            ...(item.isCritic && isCriticStyle)
                        }}

                        onClick={() => {
                            popover.handleClose();
                            item.handler?.();
                        }}
                    >

                        {item.icon && <ListItemIcon
                            sx={{
                                ...(item.isCritic && isCriticStyle)
                            }}
                        >{item.icon}</ListItemIcon>}
                        {item.label}

                    </MenuItem>)
                })}
            </Menu>


        </>
    );
};

ActionsButton.propTypes = {
    actions: PropTypes.array,
    label: PropTypes.string
};
