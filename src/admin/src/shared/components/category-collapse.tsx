import ChevronDownIcon from "@heroicons/react/24/outline/ChevronDownIcon";
import ChevronRightIcon from "@heroicons/react/24/outline/ChevronRightIcon";
import { Box, ButtonBase, Chip, Collapse, Stack, SvgIcon, Typography } from "@mui/material"
import { ReactNode, useCallback, useState } from "react";

interface CategoryCollapseProps {
    openImmediately?: boolean;
    label: ReactNode,
    children: React.ReactNode
}
export const CategoryCollapse = (props: CategoryCollapseProps) => {

    const { openImmediately = false, label, children } = props;
    const [open, setOpen] = useState<boolean>(openImmediately);
    const handleToggle = useCallback(
        (): void => {
            setOpen((prevOpen) => !prevOpen);
        },
        []
    );

    return (
        <Box sx={{ px: 2 }}>
            <ButtonBase
                onClick={handleToggle}
                sx={{
                    alignItems: 'center',
                    borderRadius: 1,
                    color: 'text.primary',
                    display: 'flex',
                    fontFamily: (theme) => theme.typography.fontFamily,
                    fontSize: 12,
                    fontWeight: 400,
                    justifyContent: 'flex-start',
                    py: '8px',
                    pr: '8px',
                    mt: '20px',

                    textAlign: 'left',
                    whiteSpace: 'nowrap',
                    width: '100%',

                    '&:focus-visible': {
                        backgroundColor: (theme) => theme.palette.action.selected
                    },
                }}
            >

                <Stack
                    direction="row"
                    alignItems="center"
                    justifyContent="space-between"
                    gap={1}
                    sx={{
                        // color: depth === 0 ? 'text.primary' : 'text.secondary',
                        width: '100%',
                        fontSize: 12,
                        mx: '10px',
                        transition: 'opacity 250ms ease-in-out',
                    }}
                >

                    <Typography variant="h6">
                        {label}
                    </Typography>

                    <SvgIcon
                        sx={{
                            color: 'neutral.500',
                            fontSize: 16,
                            transition: 'opacity 250ms ease-in-out'
                        }}
                    >
                        {open ? <ChevronDownIcon /> : <ChevronRightIcon />}
                    </SvgIcon>
                </Stack>

            </ButtonBase>
            <Collapse
                in={open}
                unmountOnExit
            >
                {children}
            </Collapse>
        </Box>
    )

}