'use client';

import type { FC, ReactNode } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import PropTypes from 'prop-types';
import { paths } from '@/lib/paths';
import { useAuth } from '@/core-features/auth/use-auth';

interface GuestGuardProps {
    children: ReactNode;
}

export const GuestGuard: FC<GuestGuardProps> = (props) => {
    const { children } = props;
    const { isAuthenticated } = useAuth();
    const router = useRouter();
    const [checked, setChecked] = useState<boolean>(false);

    const check = useCallback(
        () => {
            if (isAuthenticated) {
                //TODO:redirect to first page
                //router.replace(paths.dashboard.index);
            } else {
                setChecked(true);
            }
        },
        [isAuthenticated, router]
    );

    // Only check on mount, this allows us to redirect the user manually when auth state changes
    useEffect(
        () => {
            check();
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        []
    );

    if (!checked) {
        return null;
    }

    // If got here, it means that the redirect did not occur, and that tells us that the user is
    // not authenticated / authorized.

    return <>{children}</>;
};

