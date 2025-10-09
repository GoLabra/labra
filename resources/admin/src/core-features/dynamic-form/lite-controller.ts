import { GenericEvent } from "@/lib/utils/event";
import { useCallback, useRef } from "react";
import { Control, useFormContext, useWatch } from "react-hook-form";


interface useLiteControllerProps {
    name: string;
    disabled?: boolean;
}
export const useLiteController = <T = any>(props: useLiteControllerProps) => {

    const { name, disabled } = props;
	const formContext = useFormContext();
	const control = formContext.control;
	
    const value = useWatch({
        control,
        name,
        disabled: disabled,
        defaultValue: control._formValues[name],
        exact: true,
    });


    var registerProps = useRef(control.register(name,{value}));

    const onChange = useCallback((event: GenericEvent<T> | any) => {
        registerProps.current.onChange(event);
    }, []);

    const onBlur = useCallback((event: any) => {
        registerProps.current.onBlur(event);
    }, []);

    return {
        value: value as T ?? control._formValues[name] as T,
        name: props.name,
        disabled,
        onChange,
        onBlur,
    }
}
