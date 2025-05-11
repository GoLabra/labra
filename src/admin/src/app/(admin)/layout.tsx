'use client'

import { ApplicationStatusSubscription } from "@/core-features/subscription/application-status-subscription";
import { UserMessagesSubscription } from "@/core-features/subscription/user-messages-subscription";
// import { SubscriptionProvider } from "@/core-features/subscription/centrifugo-provider";
import { SubscriptionProvider } from "@/core-features/subscription/graphql-provider";


interface RootLayoutProps {
    children: React.ReactNode;
}
export default function RootLayout({ children }: RootLayoutProps) {

    return (
        <SubscriptionProvider>
            <ApplicationStatusSubscription />
            <UserMessagesSubscription />
            {children}
        </SubscriptionProvider>
    );
}
