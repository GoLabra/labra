import { PropsWithChildren, ReactNode, useMemo, useState } from "react";
import { useMyDialogContext } from "./use-my-dialog-context";
import { Box, ClickAwayListener, DialogTitle, Divider, IconButton, Stack, Tooltip } from "@mui/material";
import CancelIcon from '@mui/icons-material/Cancel';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';

interface DynamicDialogHeaderProps {
    whatsThis?: string;
}
export const DynamicDialogHeader = (props: PropsWithChildren<DynamicDialogHeaderProps>) => {

    const [whatsThisOpen, setWhatsThisOpen] = useState(false);
    const myDialogContext = useMyDialogContext();

    return (
        <>
            <DialogTitle
                m={0}
                p={0}>

                <Stack
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center">

                    <Box>
                        {props.children}
                    </Box>

                    <Stack direction="row" gap={2} alignItems="center">

                        {props.whatsThis && (<ClickAwayListener onClickAway={() => setWhatsThisOpen(false)}>

                            <Tooltip
                                onClose={() => setWhatsThisOpen(false)}
                                open={whatsThisOpen}
                                disableFocusListener
                                disableHoverListener
                                disableTouchListener
                                title={props.whatsThis}
                                slotProps={{
                                    popper: {
                                        disablePortal: true,
                                    },
                                }}
                            >
                                <IconButton
                                    aria-label="close"
                                    onClick={() => setWhatsThisOpen(true)}
                                    color="secondary"
                                    sx={{
                                        padding: 0
                                    }}
                                >
                                    <HelpOutlineIcon />
                                </IconButton>
                            </Tooltip>

                        </ClickAwayListener>)}

                        <IconButton
                            aria-label="close"
                            onClick={() => myDialogContext.close()}
                            sx={{
                                padding: 0
                            }}
                        >
                            <CancelIcon />
                        </IconButton>
                    </Stack>
                </Stack>
            </DialogTitle>
            <Divider />
        </>
    )
}
