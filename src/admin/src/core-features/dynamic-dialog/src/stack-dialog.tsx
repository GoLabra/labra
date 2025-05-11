import React, { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Slide, styled } from '@mui/material';

const SlideInDialog = styled(Dialog)(({ theme }) => ({
    transform: 'translateY(-100%)',
    animation: 'slideIn 500ms forwards',
    '@keyframes slideIn': {
        '0%': {
            transform: 'translateY(-100%)',
        },
        '100%': {
            transform: 'translateY(0)',
        },
    },
}));

export const StackDialog = () => {
    //const classes = styles;
    const [openFirstDialog, setOpenFirstDialog] = useState(true);
    const [openSecondDialog, setOpenSecondDialog] = useState(false);

    const handleOpenSecondDialog = () => {
        setOpenSecondDialog(true);
        setTimeout(() => {
            setOpenFirstDialog(false);
        }, 500); // Match the duration of the slide-in animation
    };

    return (
        <>
            <Button onClick={() => setOpenFirstDialog(true)} >test</Button>

            <Dialog
                open={openFirstDialog}
                onClose={() => setOpenFirstDialog(false)}
                keepMounted
            >
                <DialogTitle>Main Dialog</DialogTitle>
                <DialogContent>
                    This is the main dialog. Click the button to open the next dialog.
                    <Button onClick={() => setOpenSecondDialog(true)} >test</Button>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleOpenSecondDialog} color="primary">
                        Next
                    </Button>
                </DialogActions>
            </Dialog>

            <SlideInDialog
                open={openSecondDialog}
                onClose={() => setOpenSecondDialog(false)}
                keepMounted
                BackdropProps={{
                    style: { zIndex: 1 }, // Ensures that both dialogs share the same backdrop
                }}
            >
                <DialogTitle>Second Dialog</DialogTitle>
                <DialogContent>
                    This dialog slides in from the top to the center.
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpenSecondDialog(false)} color="primary">
                        Close
                    </Button>
                </DialogActions>
            </SlideInDialog>
        </>
    );
};