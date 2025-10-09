'use client'

import { Centrifuge, ConnectedContext, DisconnectedContext, ErrorContext, State } from 'centrifuge';
import { CENTRIFUGO_URL } from '@/config/CONST';
import { addNotification } from '../notifications/store';
import { makeVar } from '@apollo/client';
import { PubSub, subscription } from '../pub-sub';


const lastCentrifugoError = makeVar<string | null>(null);

export class CentrifugoClient {
    private _centrifuge?: Centrifuge;
    private _onError = new PubSub<string>();

    constructor(token: string) {

        if(!CENTRIFUGO_URL){
            console.log('CENTRIFUGO_URL is not set');
            return;
        }

        this._centrifuge = new Centrifuge(CENTRIFUGO_URL, {
            token //: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM3MjIiLCJleHAiOjE3Njg3NzM0NDh9.QtKStp21ze4mX6WRAwMvw_XBlsRuRbsdUJ_YCLpr6tg'
        });

        this._centrifuge.on('connected', (param: ConnectedContext) => {
            lastCentrifugoError('');
        });

        this._centrifuge.on('error', (param: ErrorContext) => {

            if(param.type == 'transport'){

                const message = `Network error: ${param.error.message} (${CENTRIFUGO_URL})`;

                if(lastCentrifugoError() == message){
                    return;
                }

                lastCentrifugoError(message);

                this._onError.publish(message);
            }
        });
    }

    onError(callback: (context: string) => void): subscription {
        return this._onError.subscribe(callback);
    }

    connect() {
        this._centrifuge?.connect();
    }

    disconnect() {
        this._centrifuge?.disconnect();
    }

    subscribe(channel: string, onMessage: (data: any) => void) {
        if(!this._centrifuge){
            return;
        }

        const subscription = this._centrifuge.newSubscription(channel);
        subscription.on('publication', (ctx) => onMessage(ctx.data));
        subscription.subscribe();

        return () => {
            subscription.unsubscribe();
            this._centrifuge!.removeSubscription(subscription);
        };
    }

    public isConnected() {
        if(!this._centrifuge){
            return false;
        }
        return this._centrifuge.state === State.Connected;
    }
}