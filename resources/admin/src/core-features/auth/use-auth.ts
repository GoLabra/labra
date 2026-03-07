'use client'

import { AuthContext as JwtAuthContext, AuthContextType as JwtAuthContextType } from './jwt-context';
import { AuthContext as CookieAuthContext, AuthContextType as CookieAuthContextType } from './cookie-context';
import { useContext } from 'react';


type AuthContextType =
    //   | AmplifyAuthContextType
    //   | Auth0AuthContextType
    //   | FirebaseAuthContextType
    | CookieAuthContextType
    | JwtAuthContextType;

export const useAuth = <T = AuthContextType>() => useContext(CookieAuthContext) as T;
