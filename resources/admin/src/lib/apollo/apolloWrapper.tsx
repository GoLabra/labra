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
import { fromError, fromPromise } from "@apollo/client/link/utils";
import {
  clearAuthTokens,
  getAccessToken,
  refreshAdminAccessToken,
} from "@/core-features/auth/token-storage";

const getCookie = (name: string): string => {
  if (typeof document === "undefined") return "";
  const cookies = document.cookie ? document.cookie.split("; ") : [];
  for (const c of cookies) {
    const [k, ...v] = c.split("=");
    if (k === name) return decodeURIComponent(v.join("="));
  }
  return "";
};

export type ApiType = "admin" | "user";
export const ADMIN_CONTEXT = { clientName: "admin" };

const getBearerToken = () => {
  const token = getAccessToken();
  return token ? `Bearer ${token}` : "";
};

const authLink = setContext((_, { headers }) => {
  return {
    headers: {
      ...headers,
      authorization: getBearerToken(),
    },
  };
});

const csrfLink = setContext((operation, { headers }) => {
  const def: any = getMainDefinition(operation.query);
  const isMutation =
    def?.kind === "OperationDefinition" && def?.operation === "mutation";

  if (!isMutation) return { headers: { ...headers } };

  // If we're using Authorization header, backend does NOT require CSRF
  const authHeader =
    (headers?.authorization as string | undefined) ??
    (headers?.Authorization as string | undefined) ??
    "";

  if (authHeader && authHeader.trim() !== "") {
    return { headers: { ...headers } };
  }

  const csrf = getCookie("csrf_token");

  return {
    headers: {
      ...headers,
      ...(csrf ? { "X-CSRF-Token": csrf } : {}),
    },
  };
});

const queryLink = () => {
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
              authorization: getBearerToken(),
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
          csrfLink.concat(authLink).concat(httpQueryLink),
        )
      : csrfLink.concat(authLink).concat(httpQueryLink);

  return queryLink;
};

const adminLink = () => {
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
              authorization: getBearerToken(),
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
          // HTTP: CSRF + Auth + Cookies
          csrfLink.concat(authLink).concat(httpEntityLink),
        )
      : // No WS: still apply CSRF + Auth + Cookies
        csrfLink.concat(authLink).concat(httpEntityLink);

  return entityLink;
};

function makeClient() {
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
          adminLink(),
          queryLink(),
        )
      : adminLink();

  // Create error link
  const errorLink = onError((message: any) => {
    const statusCode =
      message?.networkError?.statusCode ?? message?.networkError?.response?.status;

    if (statusCode === 401) {
      const context = message?.operation?.getContext?.() ?? {};
      const authHeader =
        context?.headers?.authorization ?? context?.headers?.Authorization ?? "";

      if (!authHeader || context?.authRefreshTried) {
        clearAuthTokens();
        return;
      }

      message.operation.setContext({ authRefreshTried: true });

      return fromPromise(refreshAdminAccessToken(GRAPHQL_API_URL)).flatMap(
        (nextAccessToken: string | null) => {
          if (!nextAccessToken) {
            clearAuthTokens();
            return fromError(message.networkError ?? new Error("Unauthorized"));
          }

          const headers = message.operation.getContext().headers ?? {};
          message.operation.setContext({
            headers: {
              ...headers,
              authorization: `Bearer ${nextAccessToken}`,
            },
          });

          return message.forward(message.operation);
        },
      );
    }

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
  return (
    <ApolloNextAppProvider makeClient={makeClient}>
      {children}
    </ApolloNextAppProvider>
  );
}
