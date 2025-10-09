"use client";

import type { FC } from 'react';
import { useCallback } from 'react';
import PropTypes from 'prop-types';
import ChevronDownIcon from '@heroicons/react/24/outline/ChevronDownIcon';
import type { ButtonProps } from '@mui/material';
import { Button, ListItemIcon, ListItemText, Menu, MenuItem, SvgIcon } from '@mui/material';
import { usePopover } from '../../hooks/use-popover';
import DeleteIcon from '@mui/icons-material/Delete';

interface BulkActionsMenuProps extends ButtonProps {
    disabled?: boolean;
    onArchive?: () => void;
    onDelete?: () => void;
    selectedCount?: number;
}

export const BulkActionsMenu: FC<BulkActionsMenuProps> = (props) => {
    const { disabled, onArchive, onDelete, selectedCount = 0, sx, ...other } = props;
    const popover = usePopover<HTMLButtonElement>();

    const handleArchive = useCallback(
        (): void => {
            popover.handleClose();
            onArchive?.();
        },
        [popover, onArchive]
    );

    const handleDelete = useCallback(
        (): void => {
            popover.handleClose();
            onDelete?.();
        },
        [popover, onDelete]
    );

    return (
        <>
            <Button
                disabled={disabled}
                onClick={popover.handleOpen}
                ref={popover.anchorRef}
                startIcon={(
                    <SvgIcon fontSize="small">
                        <ChevronDownIcon />
                    </SvgIcon>
                )}
                variant="outlined"
                sx={{
                    flexShrink: 0,
                    whiteSpace: 'nowrap',
                    ...sx
                }}
                {...other}
            >
                Bulk Actions
            </Button>
            <Menu
                anchorEl={popover.anchorRef.current}
                anchorOrigin={{
                    horizontal: 'right',
                    vertical: 'bottom'
                }}
                MenuListProps={{
                    dense: true
                }}
                onClose={popover.handleClose}
                open={popover.open}
                transformOrigin={{
                    horizontal: 'right',
                    vertical: 'top'
                }}
            >
                <MenuItem
                    onClick={handleDelete}
                    sx={{
                        color: "error.main"
                    }}>
                    <ListItemIcon>
                        <DeleteIcon fontSize="small" color="error" />
                    </ListItemIcon>
                    <ListItemText>
                        Delete ({selectedCount})
                    </ListItemText>
                </MenuItem>
            </Menu>
        </>
    );
};

BulkActionsMenu.propTypes = {
    disabled: PropTypes.bool,
    onArchive: PropTypes.func,
    onDelete: PropTypes.func,
    selectedCount: PropTypes.number
};
