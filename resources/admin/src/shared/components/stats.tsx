import { Grid, Paper, Typography } from "@mui/material";
import { useCallback } from "react";
import { alpha } from '@mui/material/styles';

export const Stats = ({ label, value }: { label: string; value: string }) => (
    <Grid item xs={4}>
        <Paper
            sx={{
                alignItems: "center",
                backgroundColor: (theme) => alpha(theme.palette.neutral[600], 0.2),
                borderRadius: 1,
                p: 2,
            }}
            elevation={0}
        >
            <Typography color="text.secondary" variant="overline">
                {label}
            </Typography>
            <Typography variant="h4">{value}</Typography>
        </Paper>
    </Grid>
);