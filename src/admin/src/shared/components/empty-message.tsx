import Empty from "@/shared/components/vendor/ant-d/empty"
import { Stack, Typography } from "@mui/material"

interface EmptyMessageProps {
    label?: string;
}

export const EmptyMessage = (props: EmptyMessageProps) => {

    const { label="No Data" } = props;

    return (
        <Stack
            alignItems="center"
            sx={{
                opacity: .5
            }}>

            <Empty />

            <Typography
                variant="subtitle2"
                component="span"
                align='center'
                color="text.secondary"
                whiteSpace="nowrap"
            >
                {label}
            </Typography>
        </Stack>
    )
}