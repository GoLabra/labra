import { Switch } from "@mui/material";

interface MenuItemToggleProps{
    value: boolean;
    valueChange: (value: boolean) => void;
}
export const MenuItemToggle= (props: MenuItemToggleProps) => {

    const { value, valueChange, ...other } = props;

    return (
        <Switch
            checked={value}
            onChange={(event) => valueChange(event.target.checked)}
            sx={{
                marginLeft: '30px',
            }} />
    )
}