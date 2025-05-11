import { useEffect, useRef } from "react";
import { ApolloClient, gql, useApolloClient } from "@apollo/client";
import { useSubscription } from "./use-subscription";
import { handleApplicationStatusFailed, handleApplicationStatusGenerating, handleApplicationStatusRestarting, handleApplicationStatusReverting, handleApplicationStatusUP } from "@/store/handlers/handleApplicationStatus";
import { AppStatus } from "@/lib/apollo/graphql.entities";
import { GraphQlSubscriptionProviderType } from "./graphql-provider";
import { ENTITY_CONTEXT } from '@/lib/apollo/apolloWrapper';

const appStatusSubscription = gql`
    subscription appStatus {
        appStatus
    }
`

export const ApplicationStatusSubscription = () => {

    const client = useSubscription<GraphQlSubscriptionProviderType>();

    const handlersRef = useRef<ApolloClient<any>>(null!);
    handlersRef.current = useApolloClient();
    
    useEffect(() => {
        const cleanup = client?.onReconnected(() => {
            handleApplicationStatusUP(handlersRef.current);
        });

        return () => { cleanup(); };
    }, [client.onReconnected]);


    useEffect(() => {
        return client.subscribe('appStatus', appStatusSubscription, ENTITY_CONTEXT, (data: AppStatus ) => {
            
            if (data === AppStatus.Generating) {
                handleApplicationStatusGenerating()
                return;
            }

            if (data === AppStatus.Reverting) {
                handleApplicationStatusReverting()
                return;
            }

            if (data === AppStatus.Restarting) {
                handleApplicationStatusRestarting()
                return;
            }

            if (data === AppStatus.Fatal) {
                handleApplicationStatusFailed()
                return;
            }

            if (data === AppStatus.Up) {
                handleApplicationStatusUP(handlersRef.current);
                return;
            }
        });
    }, [client.subscribe]);

    return null;
}
