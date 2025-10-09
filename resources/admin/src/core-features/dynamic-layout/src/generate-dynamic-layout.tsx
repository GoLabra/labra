import { Box, Stack } from "@mui/material";
import React from "react";
import { dynamicLayoutItem } from "./dynamic-layout";

export const renderDynamicLayout = <T extends any>(layoutItem: dynamicLayoutItem<T> | undefined, render: (data: T) => React.ReactNode) => {
    if(!layoutItem){
        return null;
    }
   return renderDynamicLayout_internal(layoutItem, render);
}

const renderDynamicLayout_internal = <T extends any>(layoutItem: dynamicLayoutItem<T>, render: (data: T) => React.ReactNode, parentGap?: number) => {

    if (layoutItem.type === 'item' || layoutItem.type === undefined) {
        return render(layoutItem.data);
    }

    if(layoutItem.type === 'group' && layoutItem.visible ===  false) {
        return null;
    }

    if (layoutItem.type == 'group') {
        return (<Stack
            key={layoutItem.key}
            sx={{
                width: '100%'
            }}
            direction={layoutItem.direction == 'horizontal' ? 'row' : 'column'}
            gap={layoutItem.gap ?? parentGap}>
            {
                layoutItem.children.map((item: dynamicLayoutItem<T>) => renderDynamicLayout_internal(item, render, layoutItem.gap))
            }

        </Stack>
        );
    }
}
 

export const collectDynamicLayoutItems = <T extends any>(layoutItem: dynamicLayoutItem<T> | undefined): T[] => {
    let items: T[] = [];

    if(!layoutItem){
        return items;
    }

    if(layoutItem.type === 'group' && layoutItem.visible ===  false) {
        return items;
    }

    if (layoutItem.type === 'item' || layoutItem.type === undefined) {
        items.push(layoutItem.data);
    } else if (layoutItem.type === 'group') {
        layoutItem.children.forEach((childItem: dynamicLayoutItem<T>) => {
            items = items.concat(collectDynamicLayoutItems(childItem));
        });
    }

    return items;
};