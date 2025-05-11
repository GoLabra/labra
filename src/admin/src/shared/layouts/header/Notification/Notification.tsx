"use client";

import { Badge, Box, Divider, IconButton, Popover, Stack, SvgIcon, Typography } from '@mui/material';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { usePopover } from '@/hooks/use-popover';
import { NotificationItem } from './notification-item';
import { useNotifications } from '@/lib/notifications/store';
import Empty from '@/shared/components/vendor/ant-d/empty';
import { EmptyMessage } from '@/shared/components/empty-message';

interface ActionProps {
    total?: number;
    onClick?: (event: React.MouseEvent<HTMLElement>) => void;
    disableTooltip?: boolean;
}

export const Notifications = ({ total, onClick, disableTooltip = false }: ActionProps) => {

    const popover = usePopover<HTMLButtonElement>();
    const notificationMessage = useNotifications();

    return (
        <>
            <Badge
                color="success"
                variant="dot"
                invisible={notificationMessage.notifications.length == 0}
            >

                <IconButton
                    color="inherit"
                    onClick={popover.handleOpen}
                    ref={popover.anchorRef}
                    aria-label="Notifications"
                    aria-haspopup="menu"
                    aria-expanded={popover.open ? 'true' : undefined} >
                    <SvgIcon
                        color="action"
                        fontSize="small">

                        <NotificationsIcon />
                    </SvgIcon>
                </IconButton>

            </Badge>

            <Popover
                anchorEl={popover.anchorRef.current}
                anchorOrigin={{
                    horizontal: 'center',
                    vertical: 'bottom'
                }}
                disableScrollLock
                onClose={popover.handleClose}
                open={popover.open}>


                <Box
                    sx={{
                        pt: 2,
                        px: 2
                    }}
                >
                    <Typography variant="h6">
                        Notifications
                    </Typography>
                </Box>
                <Stack
                    divider={<Divider sx={{ margin: '0 20px' }} />}
                    sx={{
                        listStyle: 'none',
                        m: 0,
                        p: 0
                    }}
                >
                    {notificationMessage.notifications.map(i => (<NotificationItem key={i.id} notification={i} />))}
                </Stack>


                {notificationMessage.notifications.length == 0 && (<Stack direction="column" padding={12} alignItems="center" gap={1}>
                        <EmptyMessage label="No Notifications" />
                    </Stack>)}
            </Popover>
        </>
    )

    // <ActionItem
    // 	title='Notifications'
    // 	icon={NotificationsIcon}
    // 	onClick={onClick}
    // 	badgeContent={total}
    // 	disableTooltip={disableTooltip}
    // />
};