
import { Button, styled } from "@mui/material";

export const ToolbarBtn = styled(Button)(({ theme }) => ({
    color: theme.palette.text.secondary,
    border: '0',
    background: 'none',
    borderRadius: '10px',
    padding: '8px',
    cursor: 'pointer',
    minWidth: '35px',
    whiteSpace: 'nowrap'
}));