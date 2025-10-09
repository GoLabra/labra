import { useEffect, useRef } from "react";
import { addNotification } from "@/lib/notifications/store";
import { useSubscription } from "./use-subscription";

export const UserMessagesSubscription = () => {

    // const client = useSubscription();

    // useEffect(() => {
    //     return client.subscribe('USER_MESSAGE', (data: UserMessageMessage) => {
    //          addNotification({
    //             message: data.message,
    //             type: data.severity ?? 'info',
    //          }, true)
    //     });
    // }, [client]);

    return null;
}
