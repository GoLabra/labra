'use client'

import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { ApplicationMood } from "@/shared/layouts/application-mood/application-mood";
import { SnackbarProvider } from '@/lib/snackbarProvider';
import { AuthProvider } from '@/core-features/auth/cookie-context';
import { ApolloWrapper } from '@/lib/apollo/apolloWrapper';

export function ClientProviders({ children }: { children: React.ReactNode }) {
    return (
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <AuthProvider>
                <ApolloWrapper>
                    <SnackbarProvider>
                        <ApplicationMood />
                        {children}
                    </SnackbarProvider>
                </ApolloWrapper>
            </AuthProvider>
        </LocalizationProvider>
    );
}
