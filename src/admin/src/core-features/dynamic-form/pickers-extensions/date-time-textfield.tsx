import { localeConfig } from "@/config/locale-config";
import { useTimer } from "@/core-features/dynamic-form/use-timer";
import { isDayjsCustomValue } from "@/lib/utils/dayjs-custom-value";
import { Stack, TextField } from "@mui/material";
import { forwardRef, Fragment, useMemo } from "react";
import dayjs from 'dayjs';

interface DateTimeTextFieldProps {
    inputFormat: (value: any) => string,
    formValue: any;
    InputProps: any;
    ownerState: any;
}
export const DateTimeTextField = forwardRef<any, DateTimeTextFieldProps>((props, ref) => {

    const { inputFormat, formValue, InputProps, ownerState, ...other } = props;
    const timer = useTimer();

    const isDateTimeNow = useMemo(() => {
        const isCustomValue = isDayjsCustomValue(formValue);
        timer.setActive(isCustomValue);
        return isCustomValue;
    }, [timer.setActive, formValue]);

    const placeholder = useMemo(() => {
        return inputFormat(dayjs(timer.currentTime));
    }, [inputFormat, timer.currentTime]);


    return (
        <TextField
            {...other}
            {...(isDateTimeNow && {
                placeholder: placeholder,
                value: "",
                disabled: true,
            })}
            InputProps={{
                ref: InputProps.ref,
                endAdornment: InputProps.endAdornment
            }}
        />
    );
});
DateTimeTextField.displayName = 'DateTimeTextField';


const mergeAdornments = (...adornments: React.ReactNode[]) => {
    const nonNullAdornments = adornments.filter((el) => el != null);
    if (nonNullAdornments.length === 0) {
        return null;
    }

    if (nonNullAdornments.length === 1) {
        return nonNullAdornments[0];
    }

    return (
        <Stack direction="row">
            {nonNullAdornments.map((adornment, index) => (
                <Fragment key={index}>{adornment}</Fragment>
            ))}
        </Stack>
    );
};
