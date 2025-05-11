import { Alert as MuiAlert, Box, Stack, styled, Typography } from "@mui/material";
import CheckIcon from '@mui/icons-material/Check';
import { NotificationMessage } from "@/lib/notifications/store";
import { useMemo } from "react";
import dayjs from "dayjs";
import { localeConfig } from "@/config/locale-config";


const Alert = styled(MuiAlert)(({ theme }) => ({
    backgroundColor: theme.palette.background.paper,
    margin: '10px'
}));
interface NotificationItemProps {
    notification: NotificationMessage;
}
export const NotificationItem = ({notification}: NotificationItemProps) => {

    const time = useMemo(() => dayjs(notification.timestamp).format(localeConfig.dateTime.displayFormat), [notification.timestamp]);
    return (<Box>
        <Alert severity={notification.type}>
            <Stack direction="column" spacing={1}>
                <Typography variant="body2">
                    {time}
                </Typography>
                {notification.message}
            </Stack>
        </Alert>
    </Box>)
}