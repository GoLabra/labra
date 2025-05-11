import { SvgIcon } from "@mui/material";
import EditIcon from "@/assets/icons/iconly/bulk/edit";
import CategoryIcon from "@/assets/icons/iconly/bulk/category";

export const paths = {
    index: '/',
    auth: {
        jwt: {
            login: '/auth/jwt/login',
            //register: '/auth/jwt/register'
        },
        signup: '/auth/superuser/signup',
        changeRole: '/auth/change-role',
        // firebase: {
        //   login: '/auth/firebase/login',
        //   register: '/auth/firebase/register'
        // },
        // auth0: {
        //   callback: '/auth/auth0/callback',
        //   login: '/auth/auth0/login'
        // },
        // amplify: {
        //   confirmRegister: '/auth/amplify/confirm-register',
        //   forgotPassword: '/auth/amplify/forgot-password',
        //   login: '/auth/amplify/login',
        //   register: '/auth/amplify/register',
        //   resetPassword: '/auth/amplify/reset-password'
        // }
    },
    contentManager: {
        index: "/content-manager",
    },
    entityTypeDesigner: {
        index: "/entity-type-designer",
    },
    developer: {
        index: "/developer",
    }
};
