import { useContext } from "react";
import { SubscriptionContext as CentrifugoContext, CentrifugoSubscriptionProviderType } from "./centrifugo-provider";
import { SubscriptionContext as GraphQlContext, GraphQlSubscriptionProviderType } from "./graphql-provider";


type SubscriptionProviderType =
    GraphQlSubscriptionProviderType
    | CentrifugoSubscriptionProviderType;

export const useSubscription = <T = SubscriptionProviderType>() => useContext(GraphQlContext) as T;
