export class PubSub<T> {
    private subscriptions: Array<(context: T) => void> = [];
    
    subscribe(callback: (context: T) => void): subscription {
        this.subscriptions.push(callback);

        return {
            unsubscribe: () => this.unsubscribe(callback)
        };
    }
    
    private unsubscribe(callback: (context: T) => void) {
        this.subscriptions = this.subscriptions.filter(i => i !== callback);
    }
    
    publish(context: T) {
        this.subscriptions.forEach(i => i(context));
    }
}

export type subscription = {
    unsubscribe: () => void;
}