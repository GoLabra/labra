// src/lib/apollo/resolvers.ts
import { makeVar } from '@apollo/client';

export const messagesVar = makeVar<any[]>([]);

export const resolvers = {
    Query: {
        messages: () => messagesVar(),
    },
    Mutation: {
        addMessage: (_: any, { message }: { message: any }) => {
            const currentMessages = messagesVar();
            messagesVar([...currentMessages, message]);
            return true;
        },
        clearMessages: () => {
            messagesVar([]);
            return true;
        },
    },
};