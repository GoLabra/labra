import { useCallback, useEffect, useMemo, useRef } from "react";

const usePrevious = (value: any, initialValue: any) => {
    const ref = useRef(initialValue);
    useEffect(() => {
        ref.current = value;
    });
    return ref.current;
};

export const useEffectDebugger = (name: string, effectHook: any, dependencies: any, dependencyNames = []) => {
    const previousDeps = usePrevious(dependencies, []);

    const changedDeps = dependencies.reduce((accum: any, dependency: any, index: number) => {
        if (dependency !== previousDeps[index]) {
            const keyName = dependencyNames[index] || index;
            return {
                ...accum,
                [keyName]: {
                    before: previousDeps[index],
                    after: dependency, 
                    same: previousDeps[index] === dependency
                }
            };
        }

        return accum;
    }, {});

    if (Object.keys(changedDeps).length) {
        console.log(`[${name}] [use-memo-debugger] `, changedDeps);
    }

    useEffect(effectHook, dependencies);
};



export const useCallbackDebugger = (name: string, effectHook: any, dependencies: any, dependencyNames = []) => {
    const previousDeps = usePrevious(dependencies, []);

    const changedDeps = dependencies.reduce((accum: any, dependency: any, index: number) => {
        if (dependency !== previousDeps[index]) {
            const keyName = dependencyNames[index] || index;
            return {
                ...accum,
                [keyName]: {
                    before: previousDeps[index],
                    after: dependency, 
                    same: previousDeps[index] === dependency
                }
            };
        }

        return accum;
    }, {});

    if (Object.keys(changedDeps).length) {
        console.log(`[${name}] [use-memo-debugger] `, changedDeps);
    }

    return useCallback(effectHook, dependencies);
};


export const useMemoDebugger = (name: string, effectHook: any, dependencies: any, dependencyNames = []) => {
    const previousDeps = usePrevious(dependencies, []);

    const changedDeps = dependencies.reduce((accum: any, dependency: any, index: number) => {
        if (dependency !== previousDeps[index]) {
            const keyName = dependencyNames[index] || index;
            return {
                ...accum,
                [keyName]: {
                    before: previousDeps[index],
                    after: dependency, 
                    same: previousDeps[index] === dependency
                }
            };
        }

        return accum;
    }, {});

    if (Object.keys(changedDeps).length) {
        console.log(`[${name}] [use-memo-debugger] `, changedDeps);
    }

    return useMemo(effectHook, dependencies);
};