// import type { FC, ReactNode } from 'react';
// import { createContext, useCallback, useEffect, useReducer } from 'react';
// import PropTypes from 'prop-types';
// import Auth from '@aws-amplify/auth';
// import { amplifyConfig } from '../../config';
// import type { User } from '../../types/user';
// import { Issuer } from '../../utils/auth';

// Auth.configure({
//   userPoolId: amplifyConfig.aws_user_pools_id,
//   userPoolWebClientId: amplifyConfig.aws_user_pools_web_client_id,
//   region: amplifyConfig.aws_cognito_region
// });

// interface State {
//   isInitialized: boolean;
//   isAuthenticated: boolean;
//   user: User | null;
// }

// enum ActionType {
//   INITIALIZE = 'INITIALIZE',
//   SIGN_IN = 'SIGN_IN',
//   SIGN_OUT = 'SIGN_OUT',
// }

// type InitializeAction = {
//   type: ActionType.INITIALIZE;
//   payload: {
//     isAuthenticated: boolean;
//     user: User | null;
//   };
// };

// type SignInAction = {
//   type: ActionType.SIGN_IN;
//   payload: {
//     user: User;
//   };
// };

// type SignOutAction = {
//   type: ActionType.SIGN_OUT;
// }

// type Action =
//   | InitializeAction
//   | SignInAction
//   | SignOutAction;

// type Handler = (state: State, action: any) => State;

// const initialState: State = {
//   isAuthenticated: false,
//   isInitialized: false,
//   user: null
// };

// const handlers: Record<ActionType, Handler> = {
//   INITIALIZE: (state: State, action: InitializeAction): State => {
//     const { isAuthenticated, user } = action.payload;

//     return {
//       ...state,
//       isAuthenticated,
//       isInitialized: true,
//       user
//     };
//   },
//   SIGN_IN: (state: State, action: SignInAction): State => {
//     const { user } = action.payload;

//     return {
//       ...state,
//       isAuthenticated: true,
//       user
//     };
//   },
//   SIGN_OUT: (state: State): State => ({
//     ...state,
//     isAuthenticated: false,
//     user: null
//   })
// };

// const reducer = (state: State, action: Action): State => (
//   handlers[action.type] ? handlers[action.type](state, action) : state
// );

// export interface AuthContextType extends State {
//   issuer: Issuer.Amplify;
//   signIn: (email: string, password: string) => Promise<void>;
//   signUp: (email: string, password: string) => Promise<void>;
//   confirmSignUp: (username: string, code: string) => Promise<void>;
//   resendSignUp: (username: string) => Promise<void>;
//   forgotPassword: (username: string) => Promise<void>;
//   forgotPasswordSubmit: (username: string, code: string, newPassword: string) => Promise<void>;
//   signOut: () => Promise<void>;
// }

// export const AuthContext = createContext<AuthContextType>({
//   ...initialState,
//   issuer: Issuer.Amplify,
//   signIn: () => Promise.resolve(),
//   signUp: () => Promise.resolve(),
//   confirmSignUp: () => Promise.resolve(),
//   resendSignUp: () => Promise.resolve(),
//   forgotPassword: () => Promise.resolve(),
//   forgotPasswordSubmit: () => Promise.resolve(),
//   signOut: () => Promise.resolve()
// });

// interface AuthProviderProps {
//   children: ReactNode;
// }

// export const AuthProvider: FC<AuthProviderProps> = (props) => {
//   const { children } = props;
//   const [state, dispatch] = useReducer(reducer, initialState);

//   const initialize = useCallback(
//     async (): Promise<void> => {
//       try {
//         const user = await Auth.currentAuthenticatedUser();

//         // Here you should extract the complete user profile to make it
//         // available in your entire app.
//         // The auth state only provides basic information.

//         dispatch({
//           type: ActionType.INITIALIZE,
//           payload: {
//             isAuthenticated: true,
//             user: {
//               id: user.attributes.sub,
//               avatar: '/assets/avatars/avatar-chen-simmons.jpg',
//               email: user.attributes.email,
//               name: 'Chen Simmons'
//             }
//           }
//         });
//       } catch (error) {
//         dispatch({
//           type: ActionType.INITIALIZE,
//           payload: {
//             isAuthenticated: false,
//             user: null
//           }
//         });
//       }
//     },
//     [dispatch]
//   );

//   useEffect(
//     () => {
//       initialize();
//     },
//     // eslint-disable-next-line react-hooks/exhaustive-deps
//     []
//   );

//   const signOut = useCallback(
//     async (): Promise<void> => {
//       await Auth.signOut();
//       dispatch({
//         type: ActionType.SIGN_OUT
//       });
//     },
//     [dispatch]
//   );

//   const signIn = useCallback(
//     async (email: string, password: string): Promise<void> => {
//       const user = await Auth.signIn(email, password);

//       if (user.challengeName) {
//         console.error(`Unable to login, because challenge "${user.challengeName}" is mandated and we did not handle this case.`);
//         return;
//       }

//       dispatch({
//         type: ActionType.SIGN_IN,
//         payload: {
//           user: {
//             id: user.attributes.sub,
//             avatar: '/assets/avatars/avatar-chen-simmons.jpg',
//             email: user.attributes.email,
//             name: 'Chen Simmons'
//           }
//         }
//       });
//     },
//     [dispatch]
//   );

//   const signUp = useCallback(
//     async (email: string, password: string): Promise<void> => {
//       await Auth.signUp({
//         username: email,
//         password,
//         attributes: { email }
//       });
//     },
//     []
//   );

//   const confirmSignUp = useCallback(
//     async (username: string, code: string): Promise<void> => {
//       await Auth.confirmSignUp(username, code);
//     },
//     []
//   );

//   const resendSignUp = useCallback(
//     async (username: string): Promise<void> => {
//       await Auth.resendSignUp(username);
//     },
//     []
//   );

//   const forgotPassword = useCallback(
//     async (username: string): Promise<void> => {
//       await Auth.forgotPassword(username);
//     },
//     []
//   );

//   const forgotPasswordSubmit = useCallback(
//     async (username: string, code: string, newPassword: string): Promise<void> => {
//       await Auth.forgotPasswordSubmit(username, code, newPassword);
//     },
//     []
//   );

//   return (
//     <AuthContext.Provider
//       value={{
//         ...state,
//         issuer: Issuer.Amplify,
//         signIn,
//         signUp,
//         confirmSignUp,
//         resendSignUp,
//         forgotPassword,
//         forgotPasswordSubmit,
//         signOut
//       }}
//     >
//       {children}
//     </AuthContext.Provider>
//   );
// };

// AuthProvider.propTypes = {
//   children: PropTypes.node.isRequired
// };

// export const AuthConsumer = AuthContext.Consumer;
