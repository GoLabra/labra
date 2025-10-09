'use client';

import { useAuth } from '@/core-features/auth/use-auth';
import { useRouter } from 'next/navigation';
import type { FC, ReactNode } from 'react';
import { useCallback, useEffect, useState } from 'react';

interface GuestGuardProps {
    children: ReactNode;
}

export const GuestGuard: FC<GuestGuardProps> = (props) => {
    const { children } = props;
    const { isAuthenticated } = useAuth();
    const router = useRouter();
    const [checked, setChecked] = useState<boolean>(false);

    const check = useCallback(() => {
		if (isAuthenticated) {
			//TODO:redirect to first page
			//router.replace(paths.dashboard.index);
		} else {
			setChecked(true);
		}
	},[isAuthenticated, router]);

    // Only check on mount, this allows us to redirect the user manually when auth state changes
    useEffect(() => {
            check();
        },
	// eslint-disable-next-line react-hooks/exhaustive-deps
	[]);

    if (!checked) {
        return null;
    }

    return <>{children}</>;
};

