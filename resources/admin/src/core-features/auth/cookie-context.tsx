'use client'

import type { FC, ReactNode } from 'react';
import { createContext, useCallback, useEffect, useReducer } from 'react';
import type { User } from '../../types/user';
import { Issuer } from '@/lib/utils/auth';
import { changeRoleQuery, signInQuery, superUserSignUpQuery } from '@/lib/apollo/queries/auth';
import { gql, useLazyQuery, useMutation } from '@apollo/client';
import { ADMIN_CONTEXT } from '@/lib/apollo/apolloWrapper';
import { getJwtSub, isJwtValid } from '@/lib/utils/jwt';

export const getMeDocument = gql`query getMe {
     me { 
          id
          name
          firstName
          lastName
          email
            roles {
                id
                name
            }
      }
  }`

export const AUTH_MODE = 'cookie';
export const COOKIE_NAME = 'jwt';


const signOutQuery = gql`
    mutation SignOut {
        signOut @rest(type: "User", method: "POST", path: "/admin/logout") {
            success
        }
    }
`;

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
    authMode: typeof AUTH_MODE;
    signIn: (data: any) => Promise<void>;
    changeRole: (role: string) => Promise<void>;
    signUp: (data: any) => Promise<void>;
    signOut: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType>({
    ...initialState,
    issuer: Issuer.JWT,
    authMode: AUTH_MODE,
    signIn: () => Promise.resolve(),
    changeRole: () => Promise.resolve(),
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
    const [signOutRequest] = useMutation<any>(signOutQuery);
    const [me] = useLazyQuery<any>(getMeDocument, {
        context: ADMIN_CONTEXT
    });

    const initialize = useCallback(
        async (): Promise<void> => {
            try {

                const meResponse = await me();
                const user = meResponse?.data?.me;

                dispatch({
                    type: ActionType.INITIALIZE,
                    payload: {
                        isAuthenticated: user != null,
                        user
                    }
                });
            } catch (err) {
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
            await signInRequest({
                variables: {
                    input: {
                        ...data
                    }
                }
            });

            const meResponse = await me();
            const { roles, ...user } = meResponse.data.me;

            dispatch({
                type: ActionType.SIGN_IN,
                payload: {
                    user
                }
            });
        }, [dispatch, me, signInRequest]
    );

    const changeRole = useCallback(
        async (role: string): Promise<void> => {
            await changeRoleRequest({
                variables: {
                    input: {
                        role,
                    }
                }
            });
        }, [changeRoleRequest]
    );

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
            try {
                await signOutRequest();
            } catch {
                // best-effort cookie/session invalidation
            }
            dispatch({ type: ActionType.SIGN_OUT });
        }, [dispatch, signOutRequest]
    );

    return (
        <AuthContext.Provider
            value={{
                ...state,
                issuer: Issuer.JWT,
                authMode: AUTH_MODE,
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
