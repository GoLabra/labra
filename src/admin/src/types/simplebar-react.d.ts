declare module 'simplebar-react' {
    import { ComponentType, HTMLAttributes } from 'react';
    
    interface SimpleBarProps extends HTMLAttributes<HTMLDivElement> {
        scrollableNodeProps?: Record<string, any>;
        forceVisible?: boolean | 'x' | 'y';
        autoHide?: boolean;
    }
    
    const SimpleBar: ComponentType<SimpleBarProps>;
    export default SimpleBar;
}