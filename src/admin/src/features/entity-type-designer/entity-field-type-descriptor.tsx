import { Paper, Stack, Typography } from "@mui/material";
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

            <Paper
                sx={{
                    paddingY: '4px',
                    paddingX: '8px',
                    borderRadius: 0.4,
                    backgroundColor: (theme) => lighten(theme.palette.neutral[400], 0.8),
                    color: 'primary.main',
                }}>

                <Typography sx={{
                    fontSize: '0px'
                }}>
                    {icon}
                </Typography>
            </Paper>

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
                        // color: 'text.secondary',
                        color: 'neutral.500',
                        fontStyle: 'italic'
                    }}>{description}</Typography>)}

            </Stack>

            {end}

        </Stack>
    )
}