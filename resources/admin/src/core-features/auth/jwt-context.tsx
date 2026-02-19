'use client'

import type { FC, ReactNode } from 'react';
import { createContext, useCallback, useEffect, useReducer } from 'react';
import PropTypes from 'prop-types';
//import { authApi } from '../../api/auth';
import type { User, UserRole } from '../../types/user';
import { Issuer } from '@/lib/utils/auth';
import { changeRoleQuery, signInQuery, superUserSignUpQuery } from '@/lib/apollo/queries/auth';
import { useMutation, useQuery, useLazyQuery, gql } from '@apollo/client';
import { isJwtValid, getJwtSub } from '@/lib/utils/jwt';
import { ADMIN_CONTEXT } from '@/lib/apollo/apolloWrapper';
import { GRAPHQL_API_URL } from '@/config/CONST';
import {
    ACCESS_TOKEN_STORAGE_KEY,
    clearAuthTokens,
    getAccessToken,
    refreshAdminAccessToken,
    setAccessToken,
    setAuthTokens,
} from './token-storage';


export const getMeDocument = gql`query getMe($where:AdminUserWhereInput!) {
    adminUsers(where:$where) {
          id
          name
          firstName
          lastName
          email
      }
  }`


export const STORAGE_KEY = ACCESS_TOKEN_STORAGE_KEY;

interface State {
    isInitialized: boolean;
    isAuthenticated: boolean;
    user: User | null;
}

enum ActionType {
    INITIALIZE = 'INITIALIZE',
    SIGN_IN = 'SIGN_IN',
    SIGN_OUT = 'SIGN_OUT'
}

type InitializeAction = {
    type: ActionType.INITIALIZE;
    payload: {
        isAuthenticated: boolean;
        user: User | null;
    };
};

type SignInAction = {
    type: ActionType.SIGN_IN;
    payload: {
        user: User;
    };
};

type SignOutAction = {
    type: ActionType.SIGN_OUT;
};

type Action =
    | InitializeAction
    | SignInAction
    | SignOutAction;

type Handler = (state: State, action: any) => State;

const initialState: State = {
    isAuthenticated: false,
    isInitialized: false,
    user: null
};

const handlers: Record<ActionType, Handler> = {
    INITIALIZE: (state: State, action: InitializeAction): State => {
        const { isAuthenticated, user } = action.payload;

        return {
            ...state,
            isAuthenticated,
            isInitialized: true,
            user
        };
    },
    SIGN_IN: (state: State, action: SignInAction): State => {
        const { user } = action.payload;

        return {
            ...state,
            isAuthenticated: true,
            user,
        };
    },

    SIGN_OUT: (state: State): State => ({
        ...state,
        isAuthenticated: false,
        user: null,
    })
};

const reducer = (state: State, action: Action): State => (
    handlers[action.type] ? handlers[action.type](state, action) : state
);

export interface AuthContextType extends State {
    issuer: Issuer.JWT;
    signIn: (email: string) => Promise<void>;
    changeRole: (role: string) => Promise<void>;
    signUp: (data: any) => Promise<void>;
    signOut: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType>({
    ...initialState,
    issuer: Issuer.JWT,
    signIn: () => Promise.resolve(),
    changeRole: (role: string) => Promise.resolve(),
    signUp: () => Promise.resolve(),
    signOut: () => Promise.resolve()
});

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: FC<AuthProviderProps> = (props) => {
    const { children } = props;
    const [state, dispatch] = useReducer(reducer, initialState);

    const [signUpRequest] = useMutation<any>(superUserSignUpQuery);
    const [signInRequest] = useMutation<any>(signInQuery);
    const [changeRoleRequest] = useMutation<any>(changeRoleQuery);
    const [me] = useLazyQuery<any>(getMeDocument, {
		context: ADMIN_CONTEXT
	});
    
    const initialize = useCallback(
        async (): Promise<void> => {
            try {
                const accessToken = getAccessToken();

                if (isJwtValid(accessToken)) {
                    const meResponse = await me({ variables: { where: { email: getJwtSub(accessToken) } } });

                    dispatch({
                        type: ActionType.INITIALIZE,
                        payload: {
                            isAuthenticated: true,
                            user: meResponse.data.adminUsers[0]
                        }
                    });

                    return;
                }

                const refreshedAccessToken = await refreshAdminAccessToken(GRAPHQL_API_URL);

                if (isJwtValid(refreshedAccessToken)) {
                    const meResponse = await me({ variables: { where: { email: getJwtSub(refreshedAccessToken) } } });

                    dispatch({
                        type: ActionType.INITIALIZE,
                        payload: {
                            isAuthenticated: true,
                            user: meResponse.data.adminUsers[0]
                        }
                    });

                    return;
                }

                dispatch({
                    type: ActionType.INITIALIZE,
                    payload: {
                        isAuthenticated: false,
                        user: null
                    }
                });
            } catch (err) {
                console.error(err);
                clearAuthTokens();
                dispatch({
                    type: ActionType.INITIALIZE,
                    payload: {
                        isAuthenticated: false,
                        user: null
                    }
                });
            }
        }, [dispatch, me]
    );

    useEffect(() => {
        initialize();
    }, [initialize]);

    const signIn = useCallback(
        async (data: any): Promise<void> => {
            const result = await signInRequest({
                variables: {
                    input: {
                        ...data
                    }
                }
            });

            const accessToken = result.data.signIn.token;
            const refreshToken = result.data.signIn.refresh_token;

            if (!accessToken || !refreshToken) {
                throw new Error('Missing token pair from /admin/login');
            }

            setAuthTokens(accessToken, refreshToken);

            const meResponse = await me({ variables: { where: { email: getJwtSub(accessToken) } } });
            const {roles, ...user} = meResponse.data.adminUsers[0];

            dispatch({
                type: ActionType.SIGN_IN,
                payload: {
                    user: user
                }
            });
        }, [dispatch, me, signInRequest]);

    const changeRole = useCallback(
        async (role: string): Promise<void> => {

            const result = await changeRoleRequest({
                variables: {
                    input: {
                        role: role,
                    }
                }
            });

            const accessToken = result.data.changeRole.token;
            if (!accessToken) {
                throw new Error('Missing token from /admin/change-session-role');
            }
            setAccessToken(accessToken);
        }, [changeRoleRequest]);

    const signUp = useCallback(
        async (data: any): Promise<void> => {

            await signUpRequest({
                variables: {
                    input: {
                        ...data
                    }
                }
            });
        },
        [signUpRequest]
    );

    const signOut = useCallback(
        async (): Promise<void> => {
            clearAuthTokens();
            dispatch({ type: ActionType.SIGN_OUT });
        }, [dispatch]);

    return (
        <AuthContext.Provider
            value={{
                ...state,
                issuer: Issuer.JWT,
                signIn,
                signUp,
                signOut,
                changeRole
            }}
        >
            {state.isInitialized && children}

        </AuthContext.Provider>
    );
};

export const AuthConsumer = AuthContext.Consumer;
