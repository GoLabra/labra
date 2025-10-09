'use client'

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { ApplicationMood } from "@/shared/layouts/application-mood/application-mood";
import { SnackbarProvider } from '@/lib/snackbarProvider';
import { AuthConsumer, AuthProvider } from '@/core-features/auth/jwt-context';

export function ClientProviders({ children }: { children: React.ReactNode }) {
    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <AuthProvider>
                <SnackbarProvider>
                    <ApplicationMood />
                    {children}
                </SnackbarProvider>
            </AuthProvider>
        </LocalizationProvider>
    );
}