import { useState, useEffect, useCallback, useMemo } from "react";

export const useTimer = () => {

    const [currentTime, setCurrentTime] = useState<Date>(new Date());
    const [active, setActive] = useState<boolean>(false);

    const updateCurrentTime = useCallback(() => {
        setCurrentTime(new Date());
    }, [setCurrentTime]);

    useEffect(() => {
        let timer: NodeJS.Timeout;
        if (active) {
            timer = setInterval(() => {
                updateCurrentTime();
            }, 1000);
        }
        return () => {
            clearInterval(timer);
        };
    }, [active, updateCurrentTime]);

    return useMemo(() => ({
        currentTime,
        setActive: setActive
    }), [currentTime, setActive]);
}
