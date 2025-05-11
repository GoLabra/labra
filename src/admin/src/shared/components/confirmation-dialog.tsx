"use client";

import type { FC, ReactNode } from 'react';
import PropTypes from 'prop-types';
import ExclamationCircleIcon from '@heroicons/react/24/outline/ExclamationCircleIcon';
import ExclamationTriangleIcon from '@heroicons/react/24/outline/ExclamationTriangleIcon';
import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    PaperProps,
    Stack,
    SvgIcon,
    Typography
} from '@mui/material';
import { Key, useKeyHandler } from './key-handler';

type Variant = 'error' | 'warning' | 'info';

const iconMap: Record<Variant, ReactNode> = {
    error: (
        <SvgIcon
            color="error"
            fontSize="large"
        >
            <ExclamationCircleIcon />
        </SvgIcon>
    ),
    warning: (
        <SvgIcon
            color="warning"
            fontSize="large"
        >
            <ExclamationTriangleIcon />
        </SvgIcon>
    ),
    info: (
        <SvgIcon
            color="info"
            fontSize="large"
        >
            <ExclamationCircleIcon />
        </SvgIcon>
    )
};

interface ConfirmationDialogProps {
    message?: ReactNode;
    onCancel?: () => void;
    onConfirm?: () => void;
    open?: boolean;
    title?: string;
    variant?: Variant;
    PaperProps?: Partial<PaperProps<React.ElementType>>;
}

export const ConfirmationDialog: FC<ConfirmationDialogProps> = (props) => {
    const {
        message = '',
        onCancel,
        onConfirm,
        open = false,
        title,
        variant = 'info',
        ...other
    } = props;

    const icon = iconMap[variant];

    const keyHandler = useKeyHandler({ keyToHandle: Key.Enter, modifiers:["ctrl"], onPressed: () => onConfirm?.() });

    return (
      
            <Dialog
                maxWidth="sm"
                fullWidth
                onClose={onCancel}
                open={open}
                ref={keyHandler.setRef}
                slotProps={{
                    
                }}
                sx={{
                    '.MuiDialog-container':  {
                        alignItems: 'start'
                    }
                }}
                {...other}
            >
                <DialogTitle>
                    <Stack
                        alignItems="center"
                        direction="row"
                        spacing={2}
                    >
                        {icon}
                        <Typography variant="inherit">
                            {title}
                        </Typography>
                    </Stack>
                </DialogTitle>
                <Divider />
                <DialogContent>
                    {message}
                </DialogContent>
                <Divider />
                <DialogActions>
                    <Button
                        color="inherit"
                        onClick={onCancel}
                    >
                        Cancel
                    </Button>
                    <Button
                        onClick={onConfirm}
                        variant="contained"
                    >
                        Confirm
                    </Button>
                </DialogActions>
            </Dialog>
     
    );
};

ConfirmationDialog.propTypes = {
    message: PropTypes.string,
    onCancel: PropTypes.func,
    onConfirm: PropTypes.func,
    open: PropTypes.bool,
    title: PropTypes.string,
    variant: PropTypes.oneOf<Variant>(['error', 'warning', 'info'])
};
