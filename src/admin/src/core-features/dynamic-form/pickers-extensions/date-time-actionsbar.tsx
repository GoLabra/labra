"use client";

import {
    Box,
    Button,
    DialogActions,
    InputAdornment,
    InputAdornmentProps,
    Stack,
    TextField,
} from "@mui/material";
import {
    PickersActionBar,
    PickersActionBarAction,
    PickersActionBarProps,
} from "@mui/x-date-pickers/PickersActionBar";

export function DateTimeActionbar(props: PickersActionBarProps & {onDynamicNow: () => void}) {
    const {
        onAccept,
        onClear,
        onCancel,
        onSetToday,
        onDynamicNow,
        actions,
        className,
    } = props;

    if (actions == null || actions.length === 0) {
        return null;
    }

    const actionsBtns = actions?.map((actionType) => {
        switch (actionType as string) {
            case "clear":
                return (
                    <Button
                        key={actionType}
                        aria-haspopup="true"
                        onClick={() => {
                            onClear();
                        }}
                    >
                        CLEAR
                    </Button>
                );
            case "cancel":
                return (
                    <Button
                        key={actionType}
                        aria-haspopup="true"
                        onClick={() => {
                            onCancel();
                        }}
                    >
                        CANCEL
                    </Button>
                );
            case "accept":
                return (
                    <Button
                        key={actionType}
                        aria-haspopup="true"
                        onClick={() => {
                            onAccept();
                        }}
                    >
                        OK
                    </Button>
                );
            case "today":
                return (
                    <Button
                        key={actionType}
                        aria-haspopup="true"
                        onClick={() => {
                            onSetToday();
                        }}
                    >
                        TODAY
                    </Button>
                );
            case "dynamic-now":
                return (
                    <Button
                        key={actionType}
                        aria-haspopup="true"
                        onClick={() => {
                            onDynamicNow();
                        }}
                    >
                        DYNAMIC NOW
                    </Button>
                );
            default:
                return null;
        }
    });

    return (
        <DialogActions className={className}>
            {actionsBtns}
        </DialogActions>
    );
}
