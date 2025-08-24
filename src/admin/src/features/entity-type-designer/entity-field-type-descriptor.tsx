import { Box, Paper, Stack, Typography } from "@mui/material";
import React, { FC, ReactNode } from "react";
import { lighten } from '@mui/material/styles';
import { GiPolarStar } from "react-icons/gi";
import { DesignerEntityStatus } from "@/types/entity";

interface EntityFieldTypeDescriptorProp {
    icon?: ReactNode;
    label?: string;
    description?: string;
    end?: ReactNode;
    designerStatus?: DesignerEntityStatus
}
export default function EntityFieldTypeDescriptor(props: EntityFieldTypeDescriptorProp) {

    const { icon, label, description, designerStatus, end } = props;

    if (!icon) {
        return null;
    }

    return (
        <Stack
            width="100%"
            direction="row"
            alignItems="center"
            gap={2}
            m={1}
            p={0}
            sx={{ margin: '4px' }}>

            <Box
                sx={{
                    paddingY: '4px',
                    paddingX: '8px',
					borderRadius: 1,
					backgroundColor: 'var(--mui-palette-background-paper)'
                }}>

                <Typography sx={{
                    fontSize: '0px'
                }}>
                    {icon}
                </Typography>
            </Box>

            <Stack alignItems="start">

                <Stack direction="row" spacing={0} sx={{
                    minWidth: '120px'
                }}>
                    <Typography fontSize="small">
                        {label}
                    </Typography>
                    {designerStatus == 'new' && (
                        <Typography title="New" sx={{ position: 'relative', top: '-8px', opacity: .5 }}>
                            <GiPolarStar size={12} />
                        </Typography>)}
                </Stack>

                {description && (<Typography
                    fontSize="small"
                    sx={{
                        color: 'neutral.500',
                        fontStyle: 'italic'
                    }}>{description}</Typography>)}

            </Stack>

            {end}

        </Stack>
    )
}