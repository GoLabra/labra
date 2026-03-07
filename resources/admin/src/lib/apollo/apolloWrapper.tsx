"use client";

import { ApolloLink, HttpLink, split } from "@apollo/client";
import {
    ApolloNextAppProvider,
    NextSSRApolloClient,
    NextSSRInMemoryCache,
    SSRMultipartLink,
} from "@apollo/experimental-nextjs-app-support/ssr";
import {
    GRAPHQL_QUERY_API_URL,
    GRAPHQL_API_URL,
    GRAPHQL_ADMIN_API_URL,
} from "@/config/CONST";
import { onError } from "@apollo/client/link/error";
import { addNotification } from "../notifications/store";
import { RestLink } from "apollo-link-rest";
import { setContext } from "@apollo/client/link/context";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient } from "graphql-ws";
import { getMainDefinition } from "@apollo/client/utilities";
import { STORAGE_KEY as JWT_STORAGE_KEY } from "@/core-features/auth/jwt-context";
import { useAuth } from "@/core-features/auth/use-auth";
import { getCookie } from "@/lib/utils/cookie";

export type ApiType = "admin" | "user";
export const ADMIN_CONTEXT = { clientName: "admin" };
type AuthMode = "jwt" | "cookie";

const getBearerToken = () => {
    const token = localStorage?.getItem(JWT_STORAGE_KEY);
    return token ? `Bearer ${token}` : "";
};

const bearerLink = setContext((_, { headers }) => {
    return {
        headers: {
            ...headers,
            authorization: getBearerToken(),
        },
    };
});

const cookieLink = setContext((_, { headers }) => {
    const csrf = getCookie("csrf_token");

    return {
        headers: {
            ...headers,
            ...(csrf ? { "X-CSRF-Token": csrf } : {}),
        },
    };
});

const queryLink = (authMode: AuthMode) => {
    const authLink = authMode === "jwt" ? bearerLink : cookieLink;

    // -- http
    const httpQueryLink = new HttpLink({
        uri: GRAPHQL_QUERY_API_URL,
        credentials: "include",
    });

    // -- ws
    const wsQueryLink =
        typeof window !== "undefined"
            ? new GraphQLWsLink(
                createClient({
                    url: GRAPHQL_QUERY_API_URL!,
                    retryAttempts: 100,
                    connectionParams: () => ({
                        ...(authMode === "jwt"
                            ? { authorization: getBearerToken() }
                            : {}),
                        timeout: 30000,
                        retrying: true,
                    }),
                    shouldRetry: (_) => true,
                    retryWait: (retries) =>
                        new Promise((resolve) =>
                            setTimeout(
                                resolve, Math.min(1000 * Math.pow(2, retries), 10000),
                            ),
                        ),
                }),
            )
            : null;

    const queryLink =
        typeof window !== "undefined" && wsQueryLink != null
            ? split(
                ({ query }) => {
                    const definition = getMainDefinition(query);
                    return (
                        definition.kind === "OperationDefinition" &&
                        definition.operation === "subscription"
                    );
                },
                wsQueryLink,
                authLink.concat(httpQueryLink),
            )
            : authLink.concat(httpQueryLink);

    return queryLink;
};

const adminLink = (authMode: AuthMode) => {
    const authLink = authMode === "jwt" ? bearerLink : cookieLink;

    // -- http
    const httpEntityLink = new HttpLink({
        uri: GRAPHQL_ADMIN_API_URL,
        credentials: "include",
    });

    // -- ws
    const wsEntityLink =
        typeof window !== "undefined"
            ? new GraphQLWsLink(
                createClient({
                    url: GRAPHQL_ADMIN_API_URL!,
                    retryAttempts: 100,
                    connectionParams: () => ({
                        ...(authMode === "jwt"
                            ? { authorization: getBearerToken() }
                            : {}),
                        timeout: 30000,
                        retrying: true,
                    }),
                    shouldRetry: (_) => true,
                    retryWait: (retries) =>
                        new Promise((resolve) =>
                            setTimeout(
                                resolve,
                                Math.min(1000 * Math.pow(2, retries), 10000),
                            ),
                        ),
                }),
            )
            : null;

    const entityLink =
        typeof window !== "undefined" && wsEntityLink != null
            ? split(
                ({ query }) => {
                    const definition = getMainDefinition(query);
                    return (
                        definition.kind === "OperationDefinition" &&
                        definition.operation === "subscription"
                    );
                },
                wsEntityLink,
                authLink.concat(httpEntityLink),
            )
            : authLink.concat(httpEntityLink);

    return entityLink;
};

function makeClient(authMode: AuthMode) {
    const authLink = authMode === "jwt" ? bearerLink : cookieLink;

    // - rest
    const restLink = new RestLink({
        uri: GRAPHQL_API_URL,
        credentials: "include",
    });

    const link =
        typeof window !== "undefined"
            ? ApolloLink.split(
                (operation) => {
                    return operation.getContext().clientName === "admin";
                },
                adminLink(authMode),
                queryLink(authMode),
            )
            : adminLink(authMode);

    // Create error link
    const errorLink = onError((message: any) => {
        if (message.networkError) {
            addNotification({
                message: `Network error: ${message.networkError.message} (${GRAPHQL_QUERY_API_URL})`,
                type: "error",
            });
        }
    });

    return new NextSSRApolloClient({
        cache: new NextSSRInMemoryCache({
            typePolicies: {
                Query: {
                    fields: {},
                },
            },
        }),
        link:
            typeof window === "undefined"
                ? ApolloLink.from([
                    errorLink,
                    new SSRMultipartLink({
                        stripDefer: true,
                    }),
                    authLink.concat(restLink),
                    link!,
                ])
                : ApolloLink.from([errorLink, authLink.concat(restLink), link!]),
    });
}

export function ApolloWrapper({ children }: React.PropsWithChildren) {
    const auth = useAuth<{ authMode?: AuthMode }>();
    const authMode: AuthMode = auth?.authMode === "jwt" ? "jwt" : "cookie";

    return (
        <ApolloNextAppProvider makeClient={() => makeClient(authMode)}>
            {children}
        </ApolloNextAppProvider>
    );
}
