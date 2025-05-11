import React, { useCallback, useEffect, useRef } from 'react';
import { Box } from '@mui/material';


export enum Key {
    Enter = 'Enter',
    Escape = 'Escape',
    Space = ' ',
    ArrowUp = 'ArrowUp',
    ArrowDown = 'ArrowDown',
    ArrowLeft = 'ArrowLeft',
    ArrowRight = 'ArrowRight',
    // Add more keys as needed
  }

  interface KeyHandlerProps {
    keyToHandle: Key;
    modifiers?: Array<"ctrl" | "shift" | "alt" | "meta">;
    onPressed: (e: React.KeyboardEvent<HTMLDivElement>) => void;
}

export const useKeyHandler = (props: KeyHandlerProps) => {
    const { keyToHandle, modifiers, onPressed } = props;
    const containerRef = useRef<HTMLDivElement | null>(null);

    // Refs to hold the latest values
    const keyToHandleRef = useRef<Key>(keyToHandle);
    const modifiersRef = useRef<Array<"ctrl" | "shift" | "alt" | "meta"> | undefined>(modifiers);
    const onPressedRef = useRef<(e: React.KeyboardEvent<HTMLDivElement>) => void>(onPressed);

    // Update refs when props change
    useEffect(() => {
        keyToHandleRef.current = keyToHandle;
        modifiersRef.current = modifiers;
        onPressedRef.current = onPressed;
    }, [keyToHandle, modifiers, onPressed]);

    const handleKeyDown = useCallback((e: KeyboardEvent) => {

        const modifiersPressed = !modifiersRef.current || modifiersRef.current.every(modifier => {
            switch (modifier.toLowerCase()) {
                case 'ctrl':
                    return e.ctrlKey;
                case 'shift':
                    return e.shiftKey;
                case 'alt':
                    return e.altKey;
                case 'meta':
                    return e.metaKey; // For Command key on Mac or Windows key
                default:
                    return false;
            }
        });

        
        if (modifiersPressed && e.key === keyToHandleRef.current) {
            //e.preventDefault();
            
            if (onPressedRef.current) {
                onPressedRef.current(e as unknown as React.KeyboardEvent<HTMLDivElement>);
            }
        }
    }, []); // Empty dependency array ensures the function reference stays the same

    const setRef = useCallback((node: HTMLDivElement | null) => {
        if (containerRef.current) {
            containerRef.current.removeEventListener('keydown', handleKeyDown);
        }

        containerRef.current = node;

        if (containerRef.current) {
            containerRef.current.addEventListener('keydown', handleKeyDown);
        }
    }, [handleKeyDown]);

    return {
        setRef
    };
};
