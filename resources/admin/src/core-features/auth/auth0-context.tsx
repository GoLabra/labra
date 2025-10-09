// import type { FC, ReactNode } from 'react';
// import { createContext, useCallback, useEffect, useReducer } from 'react';
// import PropTypes from 'prop-types';
// import { Auth0Client } from '@auth0/auth0-spa-js';
// import { auth0Config } from '../../config';
// import { paths } from '../../paths';
// import type { User } from '../../types/user';
// import { Issuer } from '../../utils/auth';

// const auth0Client: Auth0Client = new Auth0Client({
//   domain: auth0Config.issuer_base_url!,
//   clientId: auth0Config.client_id!,
//   cacheLocation: 'localstorage',
//   authorizationParams: {
//     redirect_uri: auth0Config.base_url + paths.auth.auth0.callback
//   }
// });

// type AppState = {
//   returnTo?: string;
// };

// interface State {
//   isInitialized: boolean;
//   isAuthenticated: boolean;
//   user: User | null;
// }

// enum ActionType {
//   INITIALIZE = 'INITIALIZE',
//   LOGIN = 'LOGIN',
//   LOGOUT = 'LOGOUT',
// }

// type InitializeAction = {
//   type: ActionType.INITIALIZE;
//   payload: {
//     isAuthenticated: boolean;
//     user: User | null;
//   };
// };

// type LoginAction = {
//   type: ActionType.LOGIN;
//   payload: {
//     user: User;
//   };
// };

// type LogoutAction = {
//   type: ActionType.LOGOUT;
// }

// type Action =
//   | InitializeAction
//   | LoginAction
//   | LogoutAction;

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
//   LOGIN: (state: State, action: LoginAction): State => {
//     const { user } = action.payload;

//     return {
//       ...state,
//       isAuthenticated: true,
//       user
//     };
//   },
//   LOGOUT: (state: State): State => ({
//     ...state,
//     isAuthenticated: false,
//     user: null
//   })
// };

// const reducer = (state: State, action: Action): State => (
//   handlers[action.type] ? handlers[action.type](state, action) : state
// );

// export interface AuthContextType extends State {
//   issuer: Issuer.Auth0;
//   loginWithRedirect: (appState?: AppState) => Promise<void>;
//   handleRedirectCallback: () => Promise<AppState | undefined>;
//   logout: () => Promise<void>;
// }

// export const AuthContext = createContext<AuthContextType>({
//   ...initialState,
//   issuer: Issuer.Auth0,
//   loginWithRedirect: () => Promise.resolve(),
//   handleRedirectCallback: () => Promise.resolve(undefined),
//   logout: () => Promise.resolve()
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
//         await auth0Client.checkSession();

//         const isAuthenticated = await auth0Client.isAuthenticated();

//         if (isAuthenticated) {
//           const user = await auth0Client.getUser();

//           // Here you should extract the complete user profile to make it
//           // available in your entire app.
//           // The auth state only provides basic information.

//           dispatch({
//             type: ActionType.INITIALIZE,
//             payload: {
//               isAuthenticated,
//               user: {
//                 id: user!.sub as string,
//                 avatar: user!.picture,
//                 email: user!.email as string,
//                 name: 'Anika Visser',
//                 plan: 'Premium'
//               }
//             }
//           });
//         } else {
//           dispatch({
//             type: ActionType.INITIALIZE,
//             payload: {
//               isAuthenticated,
//               user: null
//             }
//           });
//         }
//       } catch (err) {
//         console.error(err);
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

//   const loginWithRedirect = useCallback(
//     async (appState?: AppState): Promise<void> => {
//       await auth0Client!.loginWithRedirect({
//         appState
//       });
//     },
//     []
//   );

//   const handleRedirectCallback = useCallback(
//     async (): Promise<AppState | undefined> => {
//       const result = await auth0Client!.handleRedirectCallback();
//       const user = await auth0Client!.getUser();

//       // Here you should extract the complete user profile to make it available in your entire app.
//       // The auth state only provides basic information.

//       dispatch({
//         type: ActionType.LOGIN,
//         payload: {
//           user: {
//             id: user!.sub as string,
//             avatar: user!.picture,
//             email: user!.email as string,
//             name: 'Anika Visser',
//             plan: 'Premium'
//           }
//         }
//       });

//       return result.appState;
//     },
//     []
//   );

//   const logout = useCallback(
//     async (): Promise<void> => {
//       await auth0Client!.logout();
//       dispatch({
//         type: ActionType.LOGOUT
//       });
//     },
//     []
//   );

//   return (
//     <AuthContext.Provider
//       value={{
//         ...state,
//         issuer: Issuer.Auth0,
//         loginWithRedirect,
//         handleRedirectCallback,
//         logout
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
