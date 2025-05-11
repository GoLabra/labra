import type { FC, ReactNode } from 'react';
import Error from 'next/error';
import PropTypes from 'prop-types';
import { Issuer } from '@/lib/utils/auth';
import { useAuth } from '@/core-features/auth/use-auth';


interface IssuerGuardProps {
    children: ReactNode;
    issuer: Issuer;
}

// This guard protects an auth page from being loaded when using a different issuer.
// For example, if we are using Auth0, and we try to load `/auth/firebase/login` path, this
// will render an error.
export const IssuerGuard: FC<IssuerGuardProps> = (props) => {
    const { children, issuer: expectedIssuer } = props;
    const auth = useAuth();

    const issuer = `${auth.issuer}`;
    
    if (expectedIssuer !== issuer) {
        return (
            <Error
                statusCode={400}
                title={`Issuer mismatch, currently using ${issuer}`}
                withDarkMode={false}
            />
        );
    }

    return <>{children}</>;
};

IssuerGuard.propTypes = {
    children: PropTypes.any,
    issuer: PropTypes.any
};
