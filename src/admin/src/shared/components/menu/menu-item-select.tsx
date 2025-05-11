import { Box, ListItemIcon, ListItemText, MenuItem, Select, Stack, Typography } from "@mui/material";
import ContentCopy from '@mui/icons-material/ContentCopy';
import { forwardRef } from "react";
import { Option } from "@/core-features/dynamic-form/form-field";

interface SettingsValueSelectProps<T> {
    value: T;
    valueChange: (value: T) => void;
    options: Array<Option<T, never>>;
}
export const MenuItemSelect = <T,>(props: SettingsValueSelectProps<T>, ref: React.Ref<HTMLLIElement>) => {
    const { value, valueChange, options, ...other } = props;

    return (
        <Select
            labelId="demo-simple-select-label"
            id="demo-simple-select"
            value={value}
            autoWidth={true}
            sx={{
                width: '90px',
                minWidth: '90px',
                marginLeft: '30px',
            }} >

            {options.map((option) => (
                <MenuItem
                    key={option.label}
                    value={option.value as string | readonly string[] | number | undefined}
                    onClick={() => valueChange(option.value)}
                    sx={{ justifyContent: 'flex-end' }}>
                    {option.label}
                </MenuItem>
            ))}

        </Select>
    )
}