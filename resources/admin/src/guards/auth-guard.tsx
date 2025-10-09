'use client';

import { useAuth } from '@/core-features/auth/use-auth';
import { paths } from '@/lib/paths';
import { Issuer } from '@/lib/utils/auth';
import { useRouter } from 'next/navigation';
import type { FC, ReactNode } from 'react';
import { useCallback, useEffect, useState } from 'react';

const loginPaths: Record<Issuer, string> = {
    //   [Issuer.Amplify]: paths.auth.amplify.login,
    //   [Issuer.Auth0]: paths.auth.auth0.login,
    //   [Issuer.Firebase]: paths.auth.firebase.login,
    [Issuer.JWT]: paths.auth.jwt.login
};

interface AuthGuardProps {
    children: ReactNode;
}
export const AuthGuard: FC<AuthGuardProps> = (props) => {
    const { children } = props;
    const { issuer, isAuthenticated } = useAuth();
    const router = useRouter();
    const [checked, setChecked] = useState<boolean>(false);

	const check = useCallback(() => {
		if (!isAuthenticated) {
			const searchParams = new URLSearchParams({ returnTo: globalThis.location.href }).toString();
			const href = loginPaths[issuer] + `?${searchParams}`;
			router.replace(href);
		} else {
			setChecked(true);
		}
	}, [isAuthenticated, issuer, router]);

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