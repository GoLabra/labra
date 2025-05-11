'use client'

import { PropsWithChildren, ReactNode, useEffect, useMemo } from "react"
import { alpha, Box, CssBaseline, LinearProgress, styled, ThemeProvider, useColorScheme } from "@mui/material";
import { useAppStatus } from "@/store/app-state/use-app-state";

const AppLoader = styled(LinearProgress)(({ theme }) => [
    {
        zIndex: 1200,
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        height: '2px',
        backgroundColor: 'transparent',
    },
    theme.applyStyles('light', {
        '& .MuiLinearProgress-bar': {
            backgroundColor: alpha(theme.palette.text.primary, 0.3),
        }
    }),
    theme.applyStyles('dark', {
        '& .MuiLinearProgress-bar': {
            backgroundColor: alpha(theme.palette.text.primary, 0.3),
        }
    }),
]);

export const ApplicationMood = () => {
    const { isCodeGenerating } = useAppStatus();
    const { mode } = useColorScheme();

    const grayScalePercentage = useMemo(()=> mode === 'dark' ? 50 : 70, [mode]);

    useEffect(() => {
        // Apply grayscale only when backendCodeGeneration is true
        document.documentElement.style.filter = isCodeGenerating ? `grayscale(${grayScalePercentage}%)` : 'none';
    
        // Cleanup function to remove grayscale
        return () => {
          document.documentElement.style.filter = 'none';
        };
    }, [isCodeGenerating, grayScalePercentage]); // Add backendCodeGeneration as dependency
      
    return <>
        {isCodeGenerating && <AppLoader color="secondary" />}
    </>
}