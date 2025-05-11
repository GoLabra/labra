"use client";

import { createDayjsCustomValue, isDayjsCustomValue } from "@/lib/utils/dayjs-custom-value";
import { Box, Divider } from "@mui/material";
import { DatePickerToolbar, DatePickerToolbarProps } from "@mui/x-date-pickers/DatePicker";
import { DateTimePickerToolbarProps } from "@mui/x-date-pickers/DateTimePicker";
import { useCallback, useMemo } from "react";
import ModeStandbyIcon from '@mui/icons-material/ModeStandby';
import dayjs from 'dayjs';
import { TimePickerToolbarProps } from "@mui/x-date-pickers/TimePicker";
import { ToggleButton } from "@/shared/components/toggle-button";

var utc = require("dayjs/plugin/utc");
dayjs.extend(utc);

interface DateTimeToolbarProps {
}
type propsType = DateTimeToolbarProps & DatePickerToolbarProps<any> & (DateTimePickerToolbarProps<any> | DatePickerToolbarProps<any> | TimePickerToolbarProps<any>)
export const DateTimeToolbar = (props: propsType) => {
    const { formValue, onChange } = props;

    const setDynamicNow = useCallback((isDynamicNow: boolean) => {
        const newValue = isDynamicNow ? createDayjsCustomValue('now()') : dayjs();
        onChange(newValue);
    }, [onChange]);

    const isDateTimeNow = useMemo(() => {
        const isCustomValue = isDayjsCustomValue(formValue);
        return isCustomValue;
    }, [formValue]);

    return (
        //  Pass the className to the root element to get correct layout
        <Box className={props.className}>
            <Box
                sx={{
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    margin: 1
                }}
            >

                <DatePickerToolbar {...props} />

                <ToggleButton
                    icon={<ModeStandbyIcon />}
                    value={isDateTimeNow}
                    onValueChanged={(value) => setDynamicNow(value)}>
                    Dynamic Now
                </ToggleButton>
            </Box>
            <Divider />
        </Box>
    );
};
