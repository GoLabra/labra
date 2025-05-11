

export type dynamicLayoutItem_Item<T> = {
    type?: 'item' | undefined;
    data: T;
}

export type dynamicLayoutItem_Group<T> = {
    type: 'group';
    direction: 'vertical' | 'horizontal';
    gap?: number;
    visible?: boolean;
    children: Array<dynamicLayoutItem<T>>
    key?: string; // used for react map
}

export type dynamicLayoutItem<T> = dynamicLayoutItem_Item<T> | dynamicLayoutItem_Group<T>;
