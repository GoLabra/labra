"use client";

import type { FC, ReactNode } from 'react';
import PropTypes from 'prop-types';
import type { ListItemButtonProps } from '@mui/material';
import { ListItemButton, ListItemIcon, ListItemText } from '@mui/material';

interface ActionListItemProps extends ListItemButtonProps {
    icon?: ReactNode;
    label: string;
}

export const ActionListItem: FC<ActionListItemProps> = (props) => {
    const { icon, label, ...other } = props;

    return (
        <ListItemButton {...other}>
            {icon && (
                <ListItemIcon>
                    {icon}
                </ListItemIcon>
            )}
            <ListItemText primary={label} />
        </ListItemButton>
    );
};

ActionListItem.propTypes = {
    label: PropTypes.string.isRequired
};
