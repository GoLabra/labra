import { Box, darken, Stack, styled, Typography } from "@mui/material"
import { dark } from "@mui/material/styles/createPalette";


const RootBox = styled(Box)(({ theme }) => ({
    borderRadius: 3,
    padding: '1px 8px',
    width: 'fit-content',
    borderWidth: 1,
    borderStyle: 'solid',
    borderColor: theme.palette.divider,
    backgroundColor: darken(theme.palette.background.paper, 0.01),
}));

interface IdProps {
    id?: string
    value: string;
    rootProps?: any;
}
export const Id = (props: IdProps) => {

    const { id, value } = props;

    return (<RootBox
        {...props.rootProps}
        sx={{
    
            ...props.rootProps?.sx
        }}>
        <Stack direction="row" gap={2} sx={{
            width: 'fit-content'
        }}>
            <Box>
                <Typography color="secondary" fontSize="small" variant="subtitle1">{id || 'ID'}</Typography>
            </Box>
            <Box>
                <Typography fontSize="small" variant="subtitle1">{value}</Typography>
            </Box>

        </Stack>
    </RootBox>)
}