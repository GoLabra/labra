import { useEffect, useRef } from 'react';

const useDidMountEffect = (func: () => void, deps: Array<any>) => {
    const didMount = useRef(false);

    useEffect(() => {
        if (didMount.current) func();
        else didMount.current = true;
    }, deps);
}

export default useDidMountEffect;