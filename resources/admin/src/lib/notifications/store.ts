// src/lib/notifications/store.ts
import { makeVar, useReactiveVar } from "@apollo/client";
import { AlertColor } from "@mui/material";
import { createId } from "@paralleldrive/cuid2";
import { useCallback } from "react";
import { enqueueSnackbar } from 'notistack';

export interface NotificationMessage {
    id: string;
    message: string;
    type: AlertColor; //'info' | 'success' | 'warning' | 'error';
    timestamp: Date;
    read?: boolean;
}

export const notificationsVar = makeVar<NotificationMessage[]>([]);

// Pure functions for state modifications
export const addNotification = (notification: Omit<NotificationMessage, 'id' | 'timestamp' | 'read'>, showSnackbar = true) => {
    const newNotification = {
        ...notification,
        id: createId(),
        timestamp: new Date(),
        read: false,
    };
    notificationsVar([newNotification, ...notificationsVar()]);
    enqueueSnackbar(newNotification.message, { variant: newNotification.type });
    return newNotification;
};

// Hook for components
export const useNotifications = () => {
    const notifications = useReactiveVar(notificationsVar);

    const markAsRead = useCallback((id: string) => {
        const updated = notificationsVar().map(n =>
            n.id === id ? { ...n, read: true } : n
        );
        notificationsVar(updated);
    }, []);

    const clearNotifications = useCallback(() => {
        notificationsVar([]);
    }, []);

    return {
        notifications,
        addNotification,
        markAsRead,
        clearNotifications
    };
};