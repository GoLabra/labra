/**
 * Copyright (c) Meta Platforms, Inc. and affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

import * as React from 'react';
import {
    ReactNode,
    useCallback,
    useEffect,
    useMemo,
    useRef,
    useState,
} from 'react';
import { ToolbarBtn } from './styles';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import { usePopover } from '@/hooks/use-popover';
import { ListItemIcon, Menu, MenuItem, SvgIcon } from '@mui/material';

type DropDownContextType = {
    close: () => void;
};

const DropDownContext = React.createContext<DropDownContextType | null>(null);

export function DropDownItem({
    onClick,
    title,
    icon
}: {
    onClick: () => void;
    title?: string;
    icon?: React.ReactNode
}) {

    const dropDownContext = React.useContext(DropDownContext);

    return (

        <MenuItem
            onClick={() => {
                onClick?.();
                dropDownContext?.close();
            }}>
            {icon && <ListItemIcon>
                {icon}
            </ListItemIcon>}
            {title}
        </MenuItem>
    );
}

export default function DropDown({
    disabled = false,
    buttonLabel,
    buttonAriaLabel,
    children,
    icon
}: {
    disabled?: boolean;
    buttonAriaLabel?: string;
    buttonLabel?: string;
    children: ReactNode;
    stopCloseOnClickSelf?: boolean;
    icon?: ReactNode;
}): JSX.Element {

    const popover = usePopover<HTMLButtonElement>();

    const onClose = useCallback(() => {
        close();
        popover.handleClose();
    }, [popover]);

    const contextValue = useMemo(() => ({
        close: onClose,
    }), [onClose]);

    return (
        <>
            <DropDownContext.Provider value={contextValue}>

                <ToolbarBtn
                    type="button"
                    disabled={disabled}
                    aria-label={buttonAriaLabel || buttonLabel}
                    onClick={popover.handleOpen}
                    ref={popover.anchorRef}
                    endIcon={<KeyboardArrowDownIcon />}>
                    {icon}
                    {buttonLabel && (
                        <span className="text dropdown-button-text">{buttonLabel}</span>
                    )}
                    <i className="chevron-down" />
                </ToolbarBtn>

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
                    {children}
                </Menu>

            </DropDownContext.Provider>
        </>
    );
}
