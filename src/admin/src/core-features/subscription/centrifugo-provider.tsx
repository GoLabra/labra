import { CentrifugoClient } from '@/lib/centrifugo/centrifugo-client';
import { addNotification } from '@/lib/notifications/store';
import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useMemo } from 'react';


export type CentrifugoSubscriptionProviderType = {
    subscribe: (channel: string, onMessage: (data: any) => void) => void;
    // isConnected: () => boolean;
    // disconnect: () => void;
}

// export const useCentrifugoSubscription = () => {
//     const client = useContext<CentrifugoSubscriptionProviderType | null>(SubscriptionContext);
//     if (!client) {
//         throw new Error('useSubscription must be used within SubscriptionProvider');
//     }
//     return client;
// };

export const SubscriptionContext = createContext<CentrifugoSubscriptionProviderType | null>(null);

interface SubscriptionProviderProps {
}
export const SubscriptionProvider = ({ children }: PropsWithChildren<SubscriptionProviderProps>) => {
    const client = useMemo(() => {
        const centrifugo = new CentrifugoClient('eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM3MjIiLCJleHAiOjE3Njg3NzM0NDh9.QtKStp21ze4mX6WRAwMvw_XBlsRuRbsdUJ_YCLpr6tg');
        centrifugo.onError((message) => {
            addNotification({
                message: message,
                type: 'error',
            }, true)
        });
        return centrifugo;
    }, []);

    useEffect(() => {
        client.connect();
        return () => client.disconnect();
    }, []);

    const subscribe = useCallback((channel: string, onMessage: (data: any) => void) => {
        return client.subscribe(channel, onMessage);
    }, []);

    // const disconnect = useCallback(() => {
    //     return client.disconnect();
    // }, []);

    // const isConnected = useCallback(() => {
    //     return client.isConnected();
    // }, []);

    const value = useMemo(() => ({
        subscribe,
        //     isConnected,
        //     disconnect
    }), [subscribe]);

    return (
        <SubscriptionContext.Provider value={value}>
            {children}
        </SubscriptionContext.Provider>
    );
};