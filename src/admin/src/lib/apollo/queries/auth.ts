import { gql } from "@apollo/client";

export const signInQuery = gql`
    fragment SignInPayload on REST {
        email: String
        password: String
    }

    mutation SignIn($input: SignInPayload!) {
        signIn( input: $input) @rest(type: "User", method: "POST", path: "/login") {
            token
        }
    }
`;



export const changeRoleQuery = gql`
    fragment ChangeRolePayload on REST {
        role: String
    }

    mutation SignIn($input: ChangeRolePayload!) {
        changeRole( input: $input) @rest(type: "User", method: "POST", path: "/change-session-role") {
            token
        }
    }
`;


export const superUserSignUpQuery = gql`
    fragment SignUpPayload on REST {
        email: String
    }

    mutation SuperuserSignUp($input: SignUpPayload!) {
        superuserSignUp( input: $input) @rest(type: "User", method: "POST", path: "/signup") {
            name
        }
    }
`;