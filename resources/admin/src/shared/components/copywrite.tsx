import { Typography } from "@mui/material"
import { PRODUCT_NAME } from '@/config/CONST';

export const Copywrite = () => {
    return (<>
        <Typography variant="body2" color="text.secondary">
                <small>
                    &copy; 2025 {PRODUCT_NAME}
                </small>
        </Typography>
    </>)
}