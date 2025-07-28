import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useMemo, useRef } from 'react';
import { ApolloLink, DocumentNode, gql, useApolloClient, useSubscription as useApolloSubscription, useSubscription } from "@apollo/client";
import { unsubscribe } from 'diagnostics_channel';
import { ADMIN_CONTEXT } from '@/lib/apollo/apolloWrapper';

export type GraphQlSubscriptionProviderType = {
    subscribe: (channel: string, query: DocumentNode, context: typeof ADMIN_CONTEXT | null, onMessage: (data: any) => void) => void;
    onReconnected: (callback: () => void) => () => void;
    // isConnected: () => boolean;
    // disconnect: () => void;
}

// Add this at the top of your file
const isWebSocketLink = (link: ApolloLink): link is ApolloLink => {
    return 'client' in link;
};

export const SubscriptionContext = createContext<GraphQlSubscriptionProviderType | null>(null);

interface SubscriptionProviderProps {
}
export const SubscriptionProvider = ({ children }: PropsWithChildren<SubscriptionProviderProps>) => {
    const apolloClient = useApolloClient();

    const reconnectCallbacks = useRef<Set<() => void>>(new Set()); 
    

    useEffect(() => {
        const link = apolloClient.link.right!.left!.left!;
        if (isWebSocketLink(link) && 'client' in link) {
            // @ts-ignore
            link.client?.on('connected', () => {
                reconnectCallbacks.current.forEach(callback => callback());
            });
            // @ts-ignore
            // link.client?.on('closed', () => {
                
            // });
        }
    }, [apolloClient]);


    const subscribe = useCallback((channel: string, query: DocumentNode, context: typeof ADMIN_CONTEXT | null, onMessage: (data: any) => void) => {

        const unsubscribe = apolloClient.subscribe({
            query: query,
            context: context ?? undefined
        }).subscribe({
            next({ data }) {
                onMessage(data[channel]);
            },
            error(err) {
                console.error('Subscription error:', err);
            }
        });

        return () => {
            unsubscribe.unsubscribe();
        };
    }, [apolloClient]);

    const onReconnected = useCallback((callback: () => void) => {
        reconnectCallbacks.current.add(callback);
        return () => {
            reconnectCallbacks.current.delete(callback);
        };
    }, [reconnectCallbacks]);

    // const disconnect = useCallback(() => {
    //      // Stop all active subscriptions
    //     apolloClient.stop();

    //     // Get the websocket link
    //     const link = apolloClient.link;
    //     if (isWebSocketLink(link) && 'client' in link) {
    //         // @ts-ignore - force close the websocket
    //         link.client?.close();
    //     }
    // }, []);

    // const isConnected = useCallback(() => {
    //     const link = apolloClient.link;
    //     if (isWebSocketLink(link) && 'client' in link) {
    //         // @ts-ignore - check websocket status
    //         return link.client?.status === 1; // WebSocket.OPEN === 1
    //     }
    //     return false;
    // }, []);

    const value = useMemo(() => ({
        subscribe,
        onReconnected,
        //     isConnected,
        //     disconnect

    }), [subscribe]);

    return (
        <SubscriptionContext.Provider value={value}>
            {children}
        </SubscriptionContext.Provider>
    );
};  